package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	locationAreas, err := getLocationAreas(url, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapB(cfg *config, args []string) error {
	if cfg.Previous == nil {
		return errors.New("you're on the first page")
	}

	locationAreas, err := getLocationAreas(*cfg.Previous, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location area name")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	locationDetail, err := getLocationAreaDetails(areaName, cfg.Cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationDetail.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := getPokemon(pokemonName, cfg.Cache)
	if err != nil {
		return err
	}

	// Calculate catch chance based on base experience
	// Higher base experience = harder to catch
	// Using a simple formula: chance decreases as base experience increases
	const maxBaseExperience = 300
	catchThreshold := maxBaseExperience - pokemon.BaseExperience

	// Random number between 0 and maxBaseExperience
	randNum := rand.Intn(maxBaseExperience)

	if randNum > catchThreshold {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)

	// Add to Pokedex
	cfg.CaughtPokemon[pokemonName] = pokemon

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]

	// Check if the pokemon has been caught
	pokemon, exists := cfg.CaughtPokemon[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Print pokemon details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details about a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}
