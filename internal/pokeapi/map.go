package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func GetLocationAreas(url string) (LocAreas, error) {
	if len(url) == 0 {
		url = locationUrl
	}
	res, err := http.Get(url)

	locs := LocAreas{}

	if err != nil {
		return locs, fmt.Errorf("%w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return locs, fmt.Errorf("%w", err)
	}

	if res.StatusCode > 299 {
		return locs, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.Unmarshal(body, &locs)

	if err != nil {
		return locs, err
	}

	return locs, nil
}
