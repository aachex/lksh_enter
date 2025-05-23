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
	players := []player{}
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
			time.Sleep(time.Minute)
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

		players = append(players, p)
		playerNames = append(playerNames, p.Name+" "+p.Surname)

		res.Body.Close()
		id++
	}

	slices.Sort(playerNames)
	for _, n := range playerNames {
		fmt.Println(n)
	}

	// fetch teams
	req := getRequest(os.Getenv("API_HOST") + "/teams")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var teams []team
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &teams)
	if err != nil {
		panic(err)
	}
	res.Body.Close()

	for _, t := range teams {
		teamId[t.Name] = t.Id
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

func getRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", os.Getenv("API_TOKEN"))
	return req
}

func getStats() {
	var teamName string
	fmt.Scan(&teamName)
	id := teamId[teamName]

	req, err := http.NewRequest(http.MethodGet, os.Getenv("API_HOST")+"/matches", nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var matches []match
	err = json.Unmarshal(b, &matches)
	if err != nil {
		panic(err)
	}

	wins := 0
	defeats := 0
	scored := 0
	missed := 0
	for _, m := range matches {
		if m.Team1Id != id && m.Team2Id != id {
			continue
		}

		if id == m.Team1Id && m.Team1Score > m.Team2Score || id == m.Team2Id && m.Team2Score > m.Team1Score {
			wins++
		} else if id == m.Team1Id && m.Team1Score < m.Team2Score || id == m.Team2Id && m.Team2Score < m.Team1Score {
			defeats++
		}

		if id == m.Team1Id {
			scored += m.Team1Score
			missed += m.Team2Score
		} else {
			scored += m.Team2Score
			missed += m.Team1Score
		}
	}

	fmt.Println(wins, defeats, scored-missed)
}
