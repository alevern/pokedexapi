package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args ...string) error
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
		"explore": {
			name:        "explore",
			description: "Lists the pokemons in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to capture a certain pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get stats about a pokemon that we captured",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Content of pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(c *config, _ ...string) error {
	if len(c.pokedex) == 0 {
		return fmt.Errorf("Pokedex is empty")
	}
	fmt.Println("Your pokedex:")
	for _, pokemon := range c.pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}

func commandInspect(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Inspect error: missing pokemon name")
	}
	pokemon, err := c.pokedex[args[0]]
	if err == false {
		return fmt.Errorf("No")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Catch error: missing pokemon name")
	}
	pokemon, err := c.apiClient.GetPokemonInfos(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	dice := rand.IntN(pokemon.BaseExperience)
	if dice <= 40 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		c.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandExplore(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Explore error: missing location area")
	}
	encounters, err := c.apiClient.ListPokemonsEncounters(args[0])
	if err != nil {
		return err
	}
	if len(encounters.PokemonEncounters) == 0 {
		fmt.Println("No pokemon present in the area")
	}
	fmt.Println("Found pokemon:")
	for _, loc := range encounters.PokemonEncounters {
		fmt.Println(" - " + loc.Pokemon.Name)
	}
	return nil
}

func commandMap(c *config, _ ...string) error {
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

func commandMapb(c *config, _ ...string) error {
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

func commandHelp(_ *config, _ ...string) error {
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

func commandExit(_ *config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
