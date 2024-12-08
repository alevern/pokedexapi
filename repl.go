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

		command, exists := Commands()[commandName]
		if exists {
			err := command.callback(c)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
