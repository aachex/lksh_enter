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

func PlayerNamesSorted(client *http.Client) ([]string, error) {
	var p Player
	playerNames := []string{}
	id := 1

loop:
	for {
		req := GetRequest(os.Getenv("API_HOST") + fmt.Sprintf("/players/%d", id))

		res, err := client.Do(req)
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
