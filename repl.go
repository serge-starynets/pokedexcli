package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"workspace/github.com/serge-starynets/pokedexcli/internal/config"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokeapi"
)

func commandExit(c *config.Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config.Config, _ []string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: List the locations areas
mapb: List the previous locations areas
explore: List all Pokemons located in an area
catch: Catch a Pokemon
inspect: Inspect a Pokemon's characteristics
pokedex: Lists all your Pokemons`)
	return nil
}

func commMap(c *config.Config, _ []string) error {
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

func commMapBack(c *config.Config, _ []string) error {
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

func commandExplore(c *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: explore <area_name>")
	}
	name := args[0]

	area, err := pokeapi.GetLocationPokemons(name, c.Cache)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Found Pokemon:")
	for _, enc := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: catch <pokemon_name>")
	}
	name := args[0]

	if _, ok := c.Pokedex[name]; ok {
		fmt.Printf("You've already caught %s\n", name)
		return nil
	}

	pokemon, err := pokeapi.GetPokemon(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	x := rand.IntN(pokemon.BaseExperience + 20)
	if x >= pokemon.BaseExperience {
		c.Pokedex[name] = pokemon
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped\n", name)
	}
	return nil
}

func commandInspect(c *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: inspect <pokemon_name>")
	}
	name := args[0]

	p, ok := c.Pokedex[name]

	if !ok {
		fmt.Printf("you have not caught %s\n", name)
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil

}

func commandPokedex(c *config.Config, _ []string) error {
	if len(c.Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range c.Pokedex {
		fmt.Printf("- %s\n", name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config.Config, args []string) error
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
	"explore": {
		name:        "explore",
		description: "List all Pokemons located in an area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Catch a Pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a Pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "poredex",
		description: "List all your Pokemons",
		callback:    commandPokedex,
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
		if err := command.callback(c, words[1:]); err != nil {
			fmt.Println(err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}
}
