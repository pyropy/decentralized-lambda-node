package lambda

import (
	"log"
	"os"
)

func parseInput() ([]byte, error) {
	input, err := os.ReadFile("/inputs/input.json")
	if err != nil {
		log.Println("Error reading input")
		return nil, err
	}

	return input, nil
}
