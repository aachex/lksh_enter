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
		req := GetRequest(os.Getenv("API_HOST") + fmt.Sprintf("/players/%d", id))

		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		switch res.StatusCode {
		case http.StatusTooManyRequests:
			res.Body.Close()
			time.Sleep(time.Minute) // слишком много запросов - временно прерываем цикл (ну почему нет эндпоинта на получение всех игроков?)
			continue
		case http.StatusNotFound:
			break loop
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, err
		}

		playerNames = append(playerNames, p.Name+" "+p.Surname)

		res.Body.Close()
		id++
	}

	slices.Sort(playerNames)
	return playerNames, nil
}

func (c *Client) MustFetch(url string, obj any) {
	req := GetRequest(url)
	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
}

func GetRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", os.Getenv("API_TOKEN"))
	return req
}
