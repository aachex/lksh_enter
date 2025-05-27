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
			wins, defeats, scored, missed := client.GetStats(client.TeamId(teamName))
			fmt.Println(wins, defeats, scored-missed)

		case "versus?":
			var id1, id2 int
			fmt.Scan(&id1, &id2)
			fmt.Println(client.Versus(client.PlayerTeam(id1), client.PlayerTeam(id2)))
		}
	}
}
