package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"workspace/github.com/serge-starynets/pokedexcli/internal/pokecache"
)

var locationUrl = "https://pokeapi.co/api/v2/location-area/"

type LocAreas struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
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
	}

	err := json.Unmarshal(body, &locs)

	if err != nil {
		return locs, err
	}

	cache.Add(url, body)

	return locs, nil
}
