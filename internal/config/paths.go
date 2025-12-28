package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Cannot determine home directory: %w", err)
	}

	dir := filepath.Join(home, ".chronicle")
	return dir, nil
}

func EnsureDataDir() (string, error) {
	dir, err := DataDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("Cannot Create Chronicle Data Directory: %w", err)
	}

	return dir, nil
}

func LogPath() (string, error) {
	dir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "chronicle.log"), nil
}
