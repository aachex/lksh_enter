package general

func Versus(team1Id, team2Id int, teams []Team, matches []Match) (cnt int) {
	for _, m := range matches {
		if m.Team1Id == team1Id && m.Team2Id == team2Id || m.Team1Id == team2Id && m.Team2Id == team1Id {
			cnt++
		}
	}
	return cnt
}
