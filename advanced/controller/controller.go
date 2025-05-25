package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aachex/lksh_enter/general"
)

type Controller struct {
	Matches    []general.Match
	TeamId     map[string]int
	PlayerTeam map[int]int
	Teams      []general.Team
	Client     http.Client
}

func (c *Controller) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats", c.GetStats)
	mux.HandleFunc("GET /versus", c.GetVersus)
	mux.HandleFunc("GET /goals", c.GetGoals)
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.TeamId[teamName]
	wins, defeats, diff := general.GetStats(teamId, c.Matches)
	w.Write(fmt.Appendf(nil, "%d %d %d", wins, defeats, diff))
}

func (c *Controller) GetVersus(w http.ResponseWriter, r *http.Request) {
	id1, err := strconv.Atoi(r.URL.Query().Get("player1_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id2, err := strconv.Atoi(r.URL.Query().Get("player2_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gamesCnt := general.Versus(id1, id2, c.Teams, c.Matches)
	w.Write(fmt.Appendf(nil, "%d", gamesCnt))
}

func (c *Controller) GetGoals(w http.ResponseWriter, r *http.Request) {
	playerId, err := strconv.Atoi(r.URL.Query().Get("player_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamId := c.PlayerTeam[playerId]

	type response struct {
		MatchId int `json:"match"`
		Time    int `json:"time"`
	}
	result := []response{}

	for _, m := range c.Matches {
		if m.Team1Id != teamId && m.Team2Id != teamId {
			continue
		}

		goals := general.MustFetch[[]general.Goal](os.Getenv("API_HOST")+fmt.Sprintf("/goals?match_id=%d", m.Id), &c.Client)
		for _, g := range goals {
			if g.PlayerId == playerId {
				result = append(result, response{m.Id, g.Minute})
			}
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
