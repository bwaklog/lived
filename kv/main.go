package main

import (
	"fmt"
)

func main() {
	store := make(kvstore.KeyValueStore)
	data := []string{"GET name", "SET name John", "GET name", "SET age 30", "GET age", "exit"}

	for _, input := range data {
		if input == "exit" {
			break
		}
		output := kvstore.ProcessInput(input, store)
		fmt.Println(output)
	}
}
