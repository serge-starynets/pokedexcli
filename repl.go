package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"workspace/github.com/serge-starynets/pokedexcli/internal/config"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokeapi"
)

func commandExit(c *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config.Config) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: List the locations areas
mapb: List the previous locations areas`)
	return nil
}

func commMap(c *config.Config) error {
	locations, err := pokeapi.GetLocationAreas(c.Next, c.Cache)
	if err != nil {
		return err
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, i := range locations.Results {
		fmt.Printf("%s\n", i.Name)
	}

	return nil
}

func commMapBack(c *config.Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := pokeapi.GetLocationAreas(c.Previous, c.Cache)
	if err != nil {
		return err
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, i := range locations.Results {
		fmt.Printf("%s\n", i.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config.Config) error
}

var commandMap = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Get help",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Get locations",
		callback:    commMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Get previous locations",
		callback:    commMapBack,
	},
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func doRepl(c *config.Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command, ok := commandMap[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := command.callback(c); err != nil {
			fmt.Println(err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}
}
