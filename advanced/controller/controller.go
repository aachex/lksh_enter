package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aachex/lksh_enter/general"
)

type Controller struct {
	Matches []general.Match
	TeamId  map[string]int
	Teams   []general.Team
	Client  http.Client
}

func (c *Controller) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats", c.GetStats)
	mux.HandleFunc("GET /versus", c.GetVersus)
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
