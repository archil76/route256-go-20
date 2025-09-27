package main

import (
	"os"

	"route256/loms/internal/app"
)

func main() {
	application, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	if err := application.ListenAndServe(); err != nil {
		panic(err)
	}
}
