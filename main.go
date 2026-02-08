package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cfg := &config{
		Cache:         pokecache.NewCache(5 * time.Minute),
		CaughtPokemon: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	}
}
