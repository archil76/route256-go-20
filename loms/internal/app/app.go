package app

import (
	"context"
	"fmt"
	"net/http"
	"route256/loms/internal/infra/pgpooler"
	"strings"
	"time"

	"net"

	desc "route256/loms/internal/api"
	"route256/loms/internal/app/server"
	orderRepository "route256/loms/internal/domain/repository/postgres/order"
	stockRepository "route256/loms/internal/domain/repository/postgres/stock"
	lomsService "route256/loms/internal/domain/service"
	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/middlewares"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config *config.Config
	server *grpc.Server
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
		),
	)

	reflection.Register(grpcServer)

	postgresDsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBMaster.User,
		c.DBMaster.Password,
		c.DBMaster.Host,
		c.DBMaster.Port,
		c.DBMaster.DBName)

	pooler, err := pgpooler.NewPooler(ctx, postgresDsn)
	if err != nil {
		return nil, fmt.Errorf("NewPool: %w", err)
	}

	newStockRepository, err := stockRepository.NewStockPostgresRepository(pooler)
	if err != nil {
		return nil,
			fmt.Errorf("NewStockPostgresRepository: %w", err)
	}

	newOrderRepository, err := orderRepository.NewOrderPostgresRepository(pooler)
	if err != nil {
		return nil,
			fmt.Errorf("NewOrderPostgresRepository: %w", err)
	}

	service := lomsService.NewLomsService(newOrderRepository, newStockRepository)

	lomsServer := server.NewServer(service)

	desc.RegisterLomsServer(grpcServer, lomsServer)

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

		gwmux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(headerMatcher),
		)

		if err1 = desc.RegisterLomsHandler(ctx, gwmux, conn); err1 != nil {
			panic(err1)
		}

		gwServer := &http.Server{
			Addr:                         fmt.Sprintf(":%s", app.config.Server.HttpPort),
			Handler:                      gwmux,
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

	fmt.Printf("app up %s:%s-%s", app.config.Server.Host, app.config.Server.HttpPort, app.config.Server.GrpcPort)
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
