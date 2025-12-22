package app

import (
	"fmt"
	"net"
	"net/http"
	"route256/cart/internal/app/server"
	cartsRepository "route256/cart/internal/domain/repository/inmemoryrepository"
	lomsservice "route256/cart/internal/domain/repository/lomsrepository"
	productservice "route256/cart/internal/domain/repository/productservicerepository"
	cartsService "route256/cart/internal/domain/service"
	"route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/middlewares"
	"route256/cart/internal/infra/http/round_trippers"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/tracer"
	"strconv"
	"time"

	"context"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	config *config.Config
	server http.Server
}

func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: c}
	app.server.Handler = app.bootstrapHandlers()

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	fmt.Printf("cart service is ready %s:%s\n", app.config.Server.Host, app.config.Server.Port)

	return app.server.Serve(l)
}

func (app *App) bootstrapHandlers() http.Handler {
	ctx := context.Background()

	transport := http.DefaultTransport
	transport = round_trippers.NewTimerRoundTripper(transport)
	transport = round_trippers.NewCounterRoundTripper(transport)
	transport = round_trippers.NewLogRoundTripper(transport)
	transport = round_trippers.NewRetryRoundTripper(transport, 3, 5*time.Second)

	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	rpsLimit, err := strconv.Atoi(app.config.ProductService.Limit)
	if err != nil {
		rpsLimit = 10
	}

	jaegerURI := fmt.Sprintf("%s:%s", app.config.Jaeger.Host, app.config.Jaeger.Port)
	t, err := tracer.NewTracer(ctx, jaegerURI)
	if err != nil {
		logger.Fatalw("can't create tracer", "err", err.Error())
	}
	defer func() {
		err := t.TracerProvider.Shutdown(ctx)
		if err != nil {
			logger.Fatalw("can't shutdown tracer", "err", err.Error())
		}
	}()

	productService := productservice.NewProductService(
		httpClient,
		app.config.ProductService.Token,
		fmt.Sprintf("%s:%s", app.config.ProductService.Host, app.config.ProductService.Port),
		rpsLimit,
	)

	lomsService, err := lomsservice.NewLomsService(
		fmt.Sprintf("%s:%s", app.config.LomsService.Host, app.config.LomsService.Port),
	)
	if err != nil {
		panic(err)
	}

	const reviewsCap = 100
	cartRepository := cartsRepository.NewCartInMemoryRepository(reviewsCap, t.Tracer)
	cartService := cartsService.NewCartsService(cartRepository, productService, lomsService, t.Tracer)

	s := server.NewServer(cartService)

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", s.AddItem)
	mux.HandleFunc("GET /user/{user_id}/cart", s.GetCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", s.DeleteItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", s.ClearCart)
	mux.HandleFunc("POST /checkout/{user_id}", s.Checkout)
	mux.HandleFunc("/debug/pprof/", s.PprofHandler)
	timerMux := middlewares.NewTimeMux(mux)
	counterMux := middlewares.NewCounterMux(timerMux)
	logMux := middlewares.NewLogMux(counterMux)

	traceMux := middlewares.NewTraceMux(logMux, t.Tracer)

	return traceMux
}
