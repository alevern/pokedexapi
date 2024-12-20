package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alevern/pokedexapi/internal/client"
)

type config struct {
	apiClient        client.Client
	nextLocationsURL *string
	prevLocationsURL *string
	pokedex          map[string]client.Pokemon
}

func startRepl(c *config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		command, exists := Commands()[commandName]
		if exists {
			err := command.callback(c, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Printf("Unknown command %v", words)
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
