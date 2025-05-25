package general

import (
	"slices"
)

func Versus(id1, id2 int, teams []Team, matches []Match) (cnt int) {
	teamId1 := 0
	teamId2 := 0
	for _, t := range teams {
		if slices.Contains(t.Players, id1) {
			teamId1 = t.Id
		}
		if slices.Contains(t.Players, id2) {
			teamId2 = t.Id
		}
	}

	for _, m := range matches {
		if m.Team1Id == teamId1 && m.Team2Id == teamId2 || m.Team1Id == teamId2 && m.Team2Id == teamId1 {
			cnt++
		}
	}
	return cnt
}
