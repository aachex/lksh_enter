package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aachex/lksh_enter/advanced/logging"
	"github.com/aachex/lksh_enter/general"
)

type Controller struct {
	client general.Client
	logger *slog.Logger
}

func New(client general.Client, logger *slog.Logger) *Controller {
	return &Controller{
		client: client,
		logger: logger,
	}
}

func (c *Controller) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats", logging.Middleware(c.GetStats, c.logger))
	mux.HandleFunc("GET /front/stats", logging.Middleware(c.GetStatsHtml, c.logger))
	mux.HandleFunc("GET /versus", logging.Middleware(c.GetVersus, c.logger))
	mux.HandleFunc("GET /front/versus", logging.Middleware(c.GetVersusHtml, c.logger))
	mux.HandleFunc("GET /goals", logging.Middleware(c.GetGoals, c.logger))
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.client.TeamId(teamName)
	wins, defeats, scored, missed := c.client.GetStats(teamId)
	w.Write(fmt.Appendf(nil, "%d %d %d", wins, defeats, scored-missed))
}

func (c *Controller) GetStatsHtml(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.client.TeamId(teamName)
	wins, defeats, scored, missed := c.client.GetStats(teamId)

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

	gamesCnt := c.client.Versus(id1, id2)
	w.Write(fmt.Appendf(nil, "%d", gamesCnt))
}

func (c *Controller) GetVersusHtml(w http.ResponseWriter, r *http.Request) {
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

	player1 := c.client.Player(id1)
	player2 := c.client.Player(id2)

	tmpl, err := template.New("versus.html").ParseFiles("advanced/html/versus.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type res struct {
		// player1 info
		Name1 string
		Team1 string

		// player2 info
		Name2 string
		Team2 string

		VersusCnt int
	}

	team1Id := c.client.PlayerTeam(id1)
	team2Id := c.client.PlayerTeam(id2)
	err = tmpl.Execute(w,
		res{
			player1.Name + " " + player1.Surname,
			c.client.Team(team1Id).Name,
			player2.Name + " " + player2.Surname,
			c.client.Team(team2Id).Name,
			c.client.Versus(team1Id, team2Id),
		})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetGoals(w http.ResponseWriter, r *http.Request) {
	playerId, err := strconv.Atoi(r.URL.Query().Get("player_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamId := c.client.PlayerTeam(playerId)

	type response struct {
		MatchId int `json:"match"`
		Time    int `json:"time"`
	}
	result := []response{}

	var matches []general.Match
	c.client.MustFetch(os.Getenv("API_HOST")+"/matches", &matches)
	for _, m := range matches {
		if m.Team1Id != teamId && m.Team2Id != teamId {
			continue
		}

		var goals []general.Goal
		c.client.MustFetch(os.Getenv("API_HOST")+fmt.Sprintf("/goals?match_id=%d", m.Id), &goals)
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
