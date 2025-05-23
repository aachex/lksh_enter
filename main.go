package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// TODO: precount teams

	var s string
	for {
		fmt.Scan(&s)
		switch s {
		case "stats?":
			getStats()
		}
	}
}

func getStats() {
	var teamName string
	fmt.Scan(&teamName)

	var teamId int
	// teamId = ...

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

	for _, m := range matches {
		if m.Team1Id != teamId && m.Team2Id != teamId {
			continue
		}
	}
}
