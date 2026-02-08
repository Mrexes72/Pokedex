package main

import "pokedexcli/internal/pokecache"

type config struct {
	Next          *string
	Previous      *string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]Pokemon
}
