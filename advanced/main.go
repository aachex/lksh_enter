package main

import (
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

	client := general.Client{}

	mux := http.NewServeMux()

	controller := controller.Controller{
		Client: client,
	}

	controller.RegisterEndpoints(mux)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + os.Getenv("SRV_PORT"),
	}

	srv.ListenAndServe()
}
