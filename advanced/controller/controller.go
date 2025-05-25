package controller

import (
	"fmt"
	"net/http"
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
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	teamName := strings.Trim(r.URL.Query().Get("team_name"), "\"")
	teamId := c.TeamId[teamName]
	wins, defeats, diff := general.GetStats(teamId, c.Matches)
	w.Write(fmt.Appendf(nil, "%d %d %d", wins, defeats, diff))
}
