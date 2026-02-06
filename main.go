package main

import (
	"fmt"
)

func main() {
	result1 := cleanInput("hello world")
	fmt.Println(result1) // [hello world]

	result2 := cleanInput("Charmander Bulbasaur PIKACHU")
	fmt.Println(result2) // [charmander bulbasaur pikachu]

	result3 := cleanInput("  test   with   spaces  ")
	fmt.Println(result3) // [test with spaces]
}
