package main

import (
	"time"
	"workspace/github.com/serge-starynets/pokedexcli/internal/config"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokecache"
)

func main() {
	conf := config.Config{
		Cache: pokecache.NewCache(5 * time.Minute),
	}
	doRepl(&conf)
}
