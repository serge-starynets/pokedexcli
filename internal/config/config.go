package config

import (
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokeapi"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
}
