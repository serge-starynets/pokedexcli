package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokecache"
)

var locationUrl = "https://pokeapi.co/api/v2/location-area/"
var pokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

type LocAreas struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type LocationPokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func fetchBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("Unknown resource")
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return body, nil
}

func GetLocationAreas(url string, cache *pokecache.Cache) (LocAreas, error) {
	if len(url) == 0 {
		url = locationUrl
	}

	locs := LocAreas{}

	body, ok := cache.Get(url)

	if !ok {
		var err error
		body, err = fetchBody(url)
		if err != nil {
			return locs, err
		}
		cache.Add(url, body)
	}

	err := json.Unmarshal(body, &locs)

	if err != nil {
		return locs, err
	}

	return locs, nil
}

func GetLocationPokemons(name string, cache *pokecache.Cache) (LocationPokemons, error) {
	pokemons := LocationPokemons{}
	locUrl := locationUrl + name

	body, ok := cache.Get(locUrl)
	if !ok {
		var err error
		body, err = fetchBody(locUrl)
		if err != nil {
			return pokemons, err
		}
		cache.Add(locUrl, body)
	}

	err := json.Unmarshal(body, &pokemons)

	if err != nil {
		return pokemons, err
	}

	return pokemons, nil

}

func GetPokemon(name string) (Pokemon, error) {
	pokemon := Pokemon{}
	url := pokemonUrl + name

	body, err := fetchBody(url)
	if err != nil {
		return pokemon, err
	}

	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
