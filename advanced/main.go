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

	// fetch teams
	var teams []general.Team
	client.MustFetch(os.Getenv("API_HOST")+"/teams", &teams)

	teamId := make(map[string]int) // get team id by name
	for _, t := range teams {
		teamId[t.Name] = t.Id
	}

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
