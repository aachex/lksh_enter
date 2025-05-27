package general

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

type Client struct {
	http.Client
}

func (c *Client) GetStats(teamId int) (int, int, int, int) {
	wins := 0
	defeats := 0
	scored := 0
	missed := 0

	matches := []Match{}
	c.MustFetch(os.Getenv("API_HOST")+"/matches", &matches)

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

	return wins, defeats, scored, missed
}

func (c *Client) Versus(team1Id, team2Id int) (cnt int) {
	var matches []Match
	c.MustFetch(os.Getenv("API_HOST")+"/matches", &matches)
	for _, m := range matches {
		if m.Team1Id == team1Id && m.Team2Id == team2Id || m.Team1Id == team2Id && m.Team2Id == team1Id {
			cnt++
		}
	}
	return cnt
}

func (c *Client) PlayerNamesSorted() ([]string, error) {
	var p Player
	playerNames := []string{}
	id := 1

loop:
	for {
		statusCode := c.MustFetch(os.Getenv("API_HOST")+fmt.Sprintf("/players/%d", id), &p)

		switch statusCode {
		case http.StatusTooManyRequests:
			time.Sleep(time.Minute) // слишком много запросов - временно прерываем цикл (ну почему нет эндпоинта на получение всех игроков?)
			continue
		case http.StatusNotFound:
			break loop
		}

		playerNames = append(playerNames, p.Name+" "+p.Surname)
		id++
	}

	slices.Sort(playerNames)
	return playerNames, nil
}

// Player returns player with given id.
func (c *Client) Player(id int) Player {
	var p Player
	c.MustFetch(os.Getenv("API_HOST")+fmt.Sprintf("/players/%d", id), &p)
	return p
}

func (c *Client) Team(id int) Team {
	var t Team
	c.MustFetch(os.Getenv("API_HOST")+fmt.Sprintf("/teams/%d", id), &t)
	return t
}

// PlayerTeam returns team.Id where team.Players contains given playerId.
func (c *Client) PlayerTeam(playerId int) int {
	var teams []Team
	c.MustFetch(os.Getenv("API_HOST")+"/teams", &teams)
	for _, t := range teams {
		if slices.Contains(t.Players, playerId) {
			return t.Id
		}
	}
	return -1 // undefined player
}

// TeamId returns team.Id where team.Name equals given teamName.
func (c *Client) TeamId(teamName string) int {
	var teams []Team
	c.MustFetch(os.Getenv("API_HOST")+"/teams", &teams)
	for _, t := range teams {
		if t.Name == teamName {
			return t.Id
		}
	}
	return -1 // undefined team
}

// MustFetch requests data from given url and parses response body to obj. Obj must be a pointer.
// It returns the status code of a response.
func (c *Client) MustFetch(url string, obj any) int {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", os.Getenv("API_TOKEN"))

	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res.StatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
	return http.StatusOK
}
