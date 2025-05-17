package utilities

import (
	"os"
	"strings"
)

func IsProduction() bool {
	env_string := strings.ToLower(os.Getenv("ENV"))
	return env_string == "prod" || env_string == "production"
}

func CreateFileIfNotExists(filename string, content []byte) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, create it with content
		return os.WriteFile(filename, content, 0644)
	}
	return nil
}
