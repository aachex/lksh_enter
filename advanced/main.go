package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aachex/lksh_enter/advanced/controller"
	"github.com/aachex/lksh_enter/general"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logOpts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &logOpts))

	mux := http.NewServeMux()

	client := general.Client{}
	controller := controller.New(client, logger)
	controller.RegisterEndpoints(mux)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + os.Getenv("SRV_PORT"),
	}

	logger.Info(fmt.Sprintf("listening: %s", os.Getenv("SRV_PORT")))
	srv.ListenAndServe()
}
