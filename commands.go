package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

func Commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays the names of the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the names of the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandMap(c *config) error {
	url := c.nextLocationsURL
	locations, err := c.apiClient.ListLocations(url)
	if err != nil {
		return err
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.prevLocationsURL = locations.Previous
	c.nextLocationsURL = locations.Next
	return nil
}

func commandMapb(c *config) error {
	url := c.prevLocationsURL
	locations, err := c.apiClient.ListLocations(url)
	if err != nil {
		return err
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.prevLocationsURL = locations.Previous
	c.nextLocationsURL = locations.Next
	return nil
}

func commandHelp(_ *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for k, v := range Commands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	fmt.Println()
	return nil
}

func commandExit(_ *config) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}
