package kvstore

import "strings"

type KeyValueStore map[string]string

// ProcessInput processes the input command and updates the KeyValueStore accordingly.
func ProcessInput(input string, store KeyValueStore) string {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "Invalid input. Please provide a valid command (e.g., GET name or SET name value)"
	}

	command := parts[0]
	key := parts[1]

	switch command {
	case "GET":
		if value, ok := store[key]; ok {
			return value
		} else {
			return "Key not found"
		}
	case "SET":
		if len(parts) < 3 {
			return "Invalid input. Please provide a value for SET command"
		}
		value := parts[2]
		store[key] = value
		return "Key set successfully"
	default:
		return "Invalid command. Supported commands are GET and SET"
	}
}
