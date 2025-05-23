package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

type player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Number  int    `json:"number"`
}

type team struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Players []int  `json:"players"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	for _, p := range players {
		fmt.Println(p.Name, p.Surname)
	}
}
