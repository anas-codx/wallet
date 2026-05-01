// Package utils
package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveJSON writes any data structure to a JSON file.
func SaveJSON(filename string, data any) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, file, 0o644)
}

// LoadJSON reads a JSON file into the provided target structure.
func LoadJSON(filename string, target any) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, target)
}

// FileExists checks if a file exists at the given path.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
