package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aachex/lksh_enter/general"
)

type Controller struct {
	TeamId     map[string]int // Get team.Id by team.Name
	PlayerTeam map[int]int    // Get team.Id by player.Id
	Client     general.Client
}

func (c *Controller) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats", c.GetStats)
	mux.HandleFunc("GET /front/stats", c.GetStatsHtml)
	mux.HandleFunc("GET /versus", c.GetVersus)
	mux.HandleFunc("GET /goals", c.GetGoals)
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.TeamId[teamName]
	wins, defeats, scored, missed := c.Client.GetStats(teamId)
	w.Write(fmt.Appendf(nil, "%d %d %d", wins, defeats, scored-missed))
}

func (c *Controller) GetStatsHtml(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.TeamId[teamName]
	wins, defeats, scored, missed := c.Client.GetStats(teamId)

	tmpl, err := template.New("stats.html").ParseFiles("advanced/html/stats.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type res struct {
		TeamName string
		Wins     int
		Defeats  int
		Scored   int
		Missed   int
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, res{teamName, wins, defeats, scored, missed})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	gamesCnt := c.Client.Versus(id1, id2)
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

	var matches []general.Match
	c.Client.MustFetch(os.Getenv("API_HOST")+"/matches", &matches)
	for _, m := range matches {
		if m.Team1Id != teamId && m.Team2Id != teamId {
			continue
		}

		var goals []general.Goal
		c.Client.MustFetch(os.Getenv("API_HOST")+fmt.Sprintf("/goals?match_id=%d", m.Id), &goals)
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
