package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/aachex/lksh_enter/general"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	client := http.Client{}

	var p general.Player
	playerNames := []string{}
	id := 1

loop:
	for {
		req := general.GetRequest(os.Getenv("API_HOST") + fmt.Sprintf("/players/%d", id))

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		switch res.StatusCode {
		case http.StatusTooManyRequests:
			res.Body.Close()
			time.Sleep(time.Minute) // слишком много запросов - временно прерываем цикл (ну почему нет эндпоинта на получение всех игроков?)
			continue
		case http.StatusNotFound:
			break loop
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, &p)
		if err != nil {
			panic(err)
		}

		playerNames = append(playerNames, p.Name+" "+p.Surname)

		res.Body.Close()
		id++
	}

	slices.Sort(playerNames)
	for _, n := range playerNames {
		fmt.Println(n)
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
			teamName = teamName[1 : len(teamName)-3]
			general.GetStats(teamName, teamId, matches)

		case "versus?":
			var id1, id2 int
			fmt.Scan(&id1, &id2)
			general.Versus(id1, id2, teams, matches)
		}
	}
}
