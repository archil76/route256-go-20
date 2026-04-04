package app

import (
	"context"
	"fmt"
	"net/http"
	commentsService "route256/comments/internal/domain/service"
	"route256/comments/internal/infra/middlewares"
	"route256/comments/internal/infra/pgpooler"
	"strings"
	"time"

	"net"

	desc "route256/comments/internal/api"
	"route256/comments/internal/app/server"
	commentsRepository "route256/comments/internal/domain/repository/postgres/comments"
	"route256/comments/internal/infra/config"
	"route256/comments/internal/infra/shard_manager"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config          *config.Config
	server          *grpc.Server
	commentsService *commentsService.CommentsService
}

func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}
	ctx := context.Background()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.Validate,
			middlewares.LogUnaryServerInterceptor,
		),
	)

	reflection.Register(grpcServer)

	shardsCount := len(c.DBShards)
	poolers := make([]shard_manager.PgPooler, shardsCount)
	for i, shardCfg := range c.DBShards {
		postgresDsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			shardCfg.User,
			shardCfg.Password,
			shardCfg.Host,
			shardCfg.Port,
			shardCfg.DBName)

		pooler, err := pgpooler.NewPooler(ctx, postgresDsn)
		if err != nil {
			return nil, fmt.Errorf("NewPool: %w", err)
		}
		poolers[i] = &pooler
	}

	shardManager := shard_manager.NewShardManager(shard_manager.GetShardIndexFromID(shardsCount), poolers)
	newCommentsRepository, err := commentsRepository.NewCommentsPostgresRepository(shardManager)
	if err != nil {
		return nil,
			fmt.Errorf("NewStockPostgresRepository: %w", err)
	}

	newCommentsService := commentsService.NewCommentsService(newCommentsRepository, shardManager)

	commentsServer := server.NewServer(newCommentsService)

	desc.RegisterCommentsServer(grpcServer, commentsServer)

	app := &App{config: c}
	app.server = grpcServer

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.GrpcPort)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		conn, err1 := grpc.NewClient(
			fmt.Sprintf(":%s", app.config.Server.GrpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err1 != nil {
			panic(err1)
		}
		ctx := context.Background()

		gwMux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(headerMatcher),
		)

		if err1 = desc.RegisterCommentsHandler(ctx, gwMux, conn); err1 != nil {
			panic(err1)
		}
		gwMux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			promhttp.Handler().ServeHTTP(w, r)
		})

		logMux := middlewares.NewLogMux(gwMux)

		gwServer := &http.Server{
			Addr:                         fmt.Sprintf(":%s", app.config.Server.HttpPort),
			Handler:                      logMux,
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  10 * time.Second,
			ReadHeaderTimeout:            10 * time.Second,
			WriteTimeout:                 10 * time.Second,
			IdleTimeout:                  10 * time.Second,
		}

		if err1 = gwServer.ListenAndServe(); err1 != nil {
			panic(err1)
		}

	}()

	fmt.Printf("loms service is ready %s:%s-%s\n", app.config.Server.Host, app.config.Server.HttpPort, app.config.Server.GrpcPort)

	if err = app.server.Serve(l); err != nil {
		return err
	}

	return nil
}

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-auth":
		return key, true
	default:
		return key, false
	}
}
