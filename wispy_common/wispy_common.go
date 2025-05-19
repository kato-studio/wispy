package wispy_common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func IsProduction() bool {
	env_string := strings.ToLower(os.Getenv("ENV"))
	return env_string == "prod" || env_string == "production"
}

func CreateFileIfNotExists(filename string, content []byte) error {
	// Check if file exists first
	if _, err := os.Stat(filename); err == nil {
		return nil // File already exists
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create file with content
	if err := os.WriteFile(filename, content, 0644); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}
