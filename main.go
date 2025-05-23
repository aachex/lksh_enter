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
		req, err := http.NewRequest(http.MethodGet, os.Getenv("API_HOST")+fmt.Sprintf("/players/%d", id), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", os.Getenv("API_TOKEN"))

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

		fmt.Println(id)
		res.Body.Close()
		id++
	}

	slices.Sort(playerNames)
	for _, n := range playerNames {
		fmt.Println(n)
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

func getStats() {
	var teamName string
	fmt.Scan(&teamName)

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
}
