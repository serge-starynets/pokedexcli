package main

import "workspace/github.com/serge-starynets/pokedexcli/internal/config"

func main() {
	conf := config.Config{}
	doRepl(&conf)
}
