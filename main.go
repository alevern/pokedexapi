package main

import (
	"time"

	"github.com/alevern/pokedexapi/internal/client"
)

func main() {
	pokeClient := client.NewClient(5 * time.Second)
	cfg := &config{
		apiClient: pokeClient,
		pokedex:   make(map[string]client.Pokemon),
	}
	startRepl(cfg)
}
