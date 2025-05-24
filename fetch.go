package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func mustFetch[T any](url string) T {
	req := getRequest(url)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var result T
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func getRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", os.Getenv("API_TOKEN"))
	return req
}
