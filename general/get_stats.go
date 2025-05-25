package general

func GetStats(teamId int, matches []Match) (int, int, int) {
	wins := 0
	defeats := 0
	scored := 0
	missed := 0
	for _, m := range matches {
		if m.Team1Id != teamId && m.Team2Id != teamId {
			continue
		}

		if teamId == m.Team1Id && m.Team1Score > m.Team2Score || teamId == m.Team2Id && m.Team2Score > m.Team1Score {
			wins++
		} else if teamId == m.Team1Id && m.Team1Score < m.Team2Score || teamId == m.Team2Id && m.Team2Score < m.Team1Score {
			defeats++
		}

		if teamId == m.Team1Id {
			scored += m.Team1Score
			missed += m.Team2Score
		} else {
			scored += m.Team2Score
			missed += m.Team1Score
		}
	}

	return wins, defeats, scored - missed
}
