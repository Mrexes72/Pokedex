package main

import "strings"

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)

	lowercased := strings.ToLower(trimmed)

	words := strings.Fields(lowercased)

	return words
}
