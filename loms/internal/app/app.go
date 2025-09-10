package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"net"

	desc "route256/loms/internal/api"
	"route256/loms/internal/app/server"
	orderRepository "route256/loms/internal/domain/repository/inmemoryrepository/order"
	stockRepository "route256/loms/internal/domain/repository/inmemoryrepository/stock"
	lomsService "route256/loms/internal/domain/service"
	"route256/loms/internal/infra/config"
	"sync/atomic"

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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.Validate,
		),
	)

	reflection.Register(grpcServer)

	var sequenceGenerator atomic.Int64
	service := lomsService.NewLomsService(orderRepository.NewOrderInMemoryRepository(100, &sequenceGenerator), stockRepository.NewStockInMemoryRepository(100))

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
		conn, err := grpc.NewClient(
			fmt.Sprintf(":%s", app.config.Server.GrpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()

		gwmux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(headerMatcher),
		)

		if err = desc.RegisterLomsHandler(ctx, gwmux, conn); err != nil {
			panic(err)
		}

		gwServer := &http.Server{
			Addr:    fmt.Sprintf(":%s", app.config.Server.HttpPort),
			Handler: gwmux,
		}

		if err = gwServer.ListenAndServe(); err != nil {
			panic(err)
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
