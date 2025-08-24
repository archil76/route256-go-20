package app

import (
	"fmt"
	"net"
	"net/http"
	"route256/cart/internal/app/server"
	cartsRepository "route256/cart/internal/domain/carts/repository"
	cartsService "route256/cart/internal/domain/carts/service"
	product_service "route256/cart/internal/domain/products/service"
	"route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/round_trippers"
	"time"
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

	fmt.Printf("app bootstrap %s:%s", app.config.Server.Host, app.config.Server.Port)

	return app.server.Serve(l)
}

func (app *App) bootstrapHandlers() http.Handler {

	transport := http.DefaultTransport
	transport = round_trippers.NewLogRoundTripper(transport)
	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := product_service.NewProductService(
		httpClient,
		app.config.ProductService.Token,
		fmt.Sprintf("%s:%s", app.config.ProductService.Host, app.config.ProductService.Port),
	)

	const reviewsCap = 100
	cartRepository := cartsRepository.NewCartInMemoryRepository(reviewsCap)
	cartService := cartsService.NewCartsService(cartRepository, productService)

	s := server.NewServer(cartService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", s.AddItem)
	mux.HandleFunc("GET /user/{user_id}/cart", s.GetCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", s.DeleteItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", s.ClearCart)

	//timerMux := middlewares.NewTimeMux(mux)
	//logMux := middlewares.NewLogMux(timerMux)

	return mux
}
