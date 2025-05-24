package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func getStats() {
	in := bufio.NewReader(os.Stdin)

	teamName, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	teamName = teamName[1 : len(teamName)-3]
	id := teamId[teamName]

	req := getRequest(os.Getenv("API_HOST") + "/matches")

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
		fmt.Println(string(b))
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
