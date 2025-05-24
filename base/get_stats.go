package main

import (
	"fmt"
)

func getStats(teamName string, teamId map[string]int, matches []match) {
	id := teamId[teamName]

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
