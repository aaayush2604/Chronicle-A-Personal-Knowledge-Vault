package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	PageSize   int  `json:"page_size"`
	ShowBanner bool `json:"show_banner"`
}

func DefaultConfig() Config {
	return Config{
		PageSize:   10,
		ShowBanner: true,
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
