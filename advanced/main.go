package main

import (
	"net/http"
	"os"

	"github.com/aachex/lksh_enter/advanced/controller"
	"github.com/aachex/lksh_enter/general"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	client := http.Client{}

	// fetch teams
	teams := general.MustFetch[[]general.Team](os.Getenv("API_HOST")+"/teams", &client)

	teamId := make(map[string]int)  // get team id by name
	playerTeam := make(map[int]int) // get team id by player id
	for _, t := range teams {
		teamId[t.Name] = t.Id
		for _, playerId := range t.Players {
			playerTeam[playerId] = t.Id
		}
	}

	// fetch matches
	matches := general.MustFetch[[]general.Match](os.Getenv("API_HOST")+"/matches", &client)

	mux := http.NewServeMux()

	controller := controller.Controller{
		Teams:      teams,
		TeamId:     teamId,
		PlayerTeam: playerTeam,
		Matches:    matches,
		Client:     client,
	}

	controller.RegisterEndpoints(mux)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + os.Getenv("SRV_PORT"),
	}

	srv.ListenAndServe()
}
