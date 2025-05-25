package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/aachex/lksh_enter/general"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	client := http.Client{}

	// print players
	players, err := general.PlayerNamesSorted(&client)
	if err != nil {
		panic(err)
	}
	for _, p := range players {
		fmt.Println(p)
	}

	// fetch teams
	teams := general.MustFetch[[]general.Team](os.Getenv("API_HOST")+"/teams", &client)
	teamId := make(map[string]int)
	for _, t := range teams {
		teamId[t.Name] = t.Id
	}

	// fetch matches
	matches := general.MustFetch[[]general.Match](os.Getenv("API_HOST")+"/matches", &client)

	var s string
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Scan(&s)
		switch s {
		case "stats?":
			teamName, err := in.ReadString('\n')
			if err != nil {
				panic(err)
			}
			teamName = teamName[1 : len(teamName)-3] // убиаем кавычки
			wins, defeats, diff := general.GetStats(teamId[teamName], matches)
			fmt.Println(wins, defeats, diff)

		case "versus?":
			var id1, id2 int
			fmt.Scan(&id1, &id2)
			fmt.Println(general.Versus(id1, id2, teams, matches))
		}
	}
}
