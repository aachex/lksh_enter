package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/joho/godotenv"
)

type player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Number  int    `json:"number"`
}

type team struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Players []int  `json:"players"`
}

type match struct {
	Id         int `json:"id"`
	Team1Id    int `json:"team1"`
	Team2Id    int `json:"team2"`
	Team1Score int `json:"team1_score"`
	Team2Score int `json:"team2_score"`
}

var client http.Client

var teamId map[string]int

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var p player
	playerNames := []string{}
	id := 1

loop:
	for {
		req := getRequest(os.Getenv("API_HOST") + fmt.Sprintf("/players/%d", id))

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
	teams := mustFetch[[]team](os.Getenv("API_HOST") + "/teams")
	teamId = make(map[string]int)
	for _, t := range teams {
		teamId[t.Name] = t.Id
		fmt.Println(t.Name, t.Id)
	}

	var s string
	for {
		fmt.Scan(&s)
		switch s {
		case "stats?":
			getStats()
		}
	}
}
