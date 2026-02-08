package main

import (
	"encoding/json"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
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

func getLocationAreas(url string, cache *pokecache.Cache) (LocationAreasResponse, error) {
	// Check if we have cached data
	if data, ok := cache.Get(url); ok {
		var locationAreas LocationAreasResponse
		err := json.Unmarshal(data, &locationAreas)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		return locationAreas, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	cache.Add(url, body)

	var locationAreas LocationAreasResponse
	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locationAreas, nil
}

func getLocationAreaDetails(areaName string, cache *pokecache.Cache) (LocationAreaDetail, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	if data, ok := cache.Get(url); ok {
		var locationDetail LocationAreaDetail
		err := json.Unmarshal(data, &locationDetail)
		if err != nil {
			return LocationAreaDetail{}, err
		}
		return locationDetail, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaDetail{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaDetail{}, err
	}

	cache.Add(url, body)

	var locationDetail LocationAreaDetail
	err = json.Unmarshal(body, &locationDetail)
	if err != nil {
		return LocationAreaDetail{}, nil
	}

	return locationDetail, nil

}

func getPokemon(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	// Check if we have cached data
	if data, ok := cache.Get(url); ok {
		var pokemon Pokemon
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// Add to cache
	cache.Add(url, body)

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
