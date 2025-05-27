package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aachex/lksh_enter/general"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client := general.Client{}

	// print players
	players, err := client.PlayerNamesSorted()
	if err != nil {
		panic(err)
	}
	for _, p := range players {
		fmt.Println(p)
	}

	// fetch teams
	var teams []general.Team
	client.MustFetch(os.Getenv("API_HOST")+"/teams", &teams)
	teamId := make(map[string]int)  // get team id by name
	playerTeam := make(map[int]int) // get team id by player id
	for _, t := range teams {
		teamId[t.Name] = t.Id
		for _, playerId := range t.Players {
			playerTeam[playerId] = t.Id
		}
	}

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
			wins, defeats, scored, missed := client.GetStats(teamId[teamName])
			fmt.Println(wins, defeats, scored-missed)

		case "versus?":
			var id1, id2 int
			fmt.Scan(&id1, &id2)
			fmt.Println(client.Versus(playerTeam[id1], playerTeam[id2]))
		}
	}
}
