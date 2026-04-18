package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	PageSize   int      `json:"page_size"`
	ShowBanner bool     `json:"show_banner"`
	Paths      []string `json:"paths"`
}

func DefaultConfig() Config {
	return Config{
		PageSize:   10,
		ShowBanner: true,
		Paths:      []string{},
	}
}

func ConfigPath() (string, error) {
	dir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	cfg := DefaultConfig()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return Config{}, fmt.Errorf("cannot open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config file: %w", err)
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// 🔥 Normalize path (core fix)
func normalizeLogPath(path string) string {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return filepath.Join(path, "chronicle.log")
	}

	// If doesn't exist but ends with slash → treat as directory
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return filepath.Join(path, "chronicle.log")
	}

	return path
}

func filterExistingPaths(paths []string) []string {
	var valid []string

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			valid = append(valid, p)
		}
	}

	return valid
}

func deduplicate(paths []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			result = append(result, p)
		}
	}
	return result
}

func SelectLogPath() (string, error) {
	defaultPath, err := DefaultLogPath()
	if err != nil {
		return "", err
	}

	cfg, err := Load()
	if err != nil {
		return "", err
	}

	// Step 1: Deduplicate
	cfg.Paths = deduplicate(cfg.Paths)

	// Step 2: Normalize ALL stored paths (🔥 fixes your bug)
	var normalizedPaths []string
	for _, p := range cfg.Paths {
		normalizedPaths = append(normalizedPaths, normalizeLogPath(p))
	}

	// Step 3: Deduplicate again (in case normalization caused duplicates)
	normalizedPaths = deduplicate(normalizedPaths)

	// Step 4: Filter only existing files
	validPaths := filterExistingPaths(normalizedPaths)

	// Step 5: Save cleaned config (self-healing)
	cfg.Paths = normalizedPaths
	if err := Save(cfg); err != nil {
		return "", err
	}

	// Step 6: Build final list (default always first)
	finalPaths := []string{defaultPath}
	for _, p := range validPaths {
		if p != defaultPath {
			finalPaths = append(finalPaths, p)
		}
	}

	fmt.Println("Select a log file:")

	for i, p := range finalPaths {
		if i == 0 {
			fmt.Printf("%d) %s (default)\n", i+1, p)
		} else {
			fmt.Printf("%d) %s\n", i+1, p)
		}
	}

	fmt.Printf("%d) Enter new path\n", len(finalPaths)+1)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Choice: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("invalid input")
	}

	// 🔥 New path
	if choice == len(finalPaths)+1 {
		fmt.Print("Enter path: ")

		newInput, _ := reader.ReadString('\n')
		newPath := strings.TrimSpace(newInput)

		if newPath == "" {
			return "", fmt.Errorf("path cannot be empty")
		}

		newPath = normalizeLogPath(newPath)

		cfg.Paths = append(cfg.Paths, newPath)
		if err := Save(cfg); err != nil {
			return "", err
		}

		return newPath, nil
	}

	// Existing selection
	if choice >= 1 && choice <= len(finalPaths) {
		selected := finalPaths[choice-1]

		// Save if not already present
		exists := false
		for _, p := range cfg.Paths {
			if p == selected {
				exists = true
				break
			}
		}

		if !exists {
			cfg.Paths = append(cfg.Paths, selected)
			if err := Save(cfg); err != nil {
				return "", err
			}
		}

		return selected, nil
	}

	return "", fmt.Errorf("invalid choice")
}

func DefaultLogPath() (string, error) {
	dir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "chronicle.log"), nil
}
