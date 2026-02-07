package app

import (
	"context"
	"fmt"
	"net/http"
	lomsService "route256/loms/internal/domain/service/loms"
	outboxService "route256/loms/internal/domain/service/outbox"
	"route256/loms/internal/infra/logger"
	"route256/loms/internal/infra/pgpooler"
	"strings"
	"time"

	"net"

	desc "route256/loms/internal/api"
	"route256/loms/internal/app/server"
	orderRepository "route256/loms/internal/domain/repository/postgres/order"
	outboxRepository "route256/loms/internal/domain/repository/postgres/outbox"
	stockRepository "route256/loms/internal/domain/repository/postgres/stock"
	"route256/loms/internal/infra/config"
	kafkaProducer "route256/loms/internal/infra/kafka/producer"
	"route256/loms/internal/infra/middlewares"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config        *config.Config
	server        *grpc.Server
	producer      *kafkaProducer.KafkaProducer
	outboxService *outboxService.OutboxService
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
			middlewares.TimerUnaryServerInterceptor,
			middlewares.CounterUnaryServerInterceptor,
			middlewares.LogUnaryServerInterceptor,
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

	producer, err := kafkaProducer.NewProducer(ctx, c.Kafka.Brokers, c.Kafka.OrderTopic)
	if err != nil {
		return nil, err
	}

	newOutboxRepository, err := outboxRepository.NewOutboxPostgresRepository(pooler)
	if err != nil {
		return nil,
			fmt.Errorf("NewOutboxPostgresRepository: %w", err)
	}

	sendinterval := 1 // можно вынести в конфиги
	newOutboxService := outboxService.NewOutboxService(ctx, newOutboxRepository, sendinterval, &producer, pooler)

	service := lomsService.NewLomsService(newOrderRepository, newStockRepository, newOutboxService, pooler)

	lomsServer := server.NewServer(service)

	newOutboxService.Start()

	desc.RegisterLomsServer(grpcServer, lomsServer)

	app := &App{config: c}
	app.server = grpcServer
	app.producer = &producer

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.GrpcPort)

	defer func() {
		if err := app.producer.Close(); err != nil {
			logger.Errorw("Error closing producer: %v", err)
		}
		logger.Infow("kafka is closed")

		app.server.GracefulStop()
		app.outboxService.Stop()

	}()

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

		if err1 = desc.RegisterLomsHandler(ctx, gwMux, conn); err1 != nil {
			panic(err1)
		}
		gwMux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			promhttp.Handler().ServeHTTP(w, r)
		})

		timerMux := middlewares.NewTimerMux(gwMux)
		counterMux := middlewares.NewCounterMux(timerMux)
		logMux := middlewares.NewLogMux(counterMux)

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
