package main

import (
	"os"
	"route256/cart/internal/infra/logger"

	"route256/cart/internal/app"
)

func main() {
	defer func() {
		_ = logger.Sync()
	}()

	application, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		logger.Fatalw("Couldn't read CONFIG_FILE", "err", err.Error())
	}

	if err = application.ListenAndServe(); err != nil {
		logger.Fatalw("Couldn't listen on :8080", "err", err.Error())
	}
}
