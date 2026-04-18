package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"chronicle/internal/config"
	"chronicle/internal/engine"
	"chronicle/internal/index"
	"chronicle/internal/store"
	"chronicle/internal/terminal"
)

func main() {
	logPath, err := config.SelectLogPath()
	fmt.Println("Log Stored in: ", logPath)
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	store, err := store.New(logPath)
	if err != nil {
		log.Fatal(err)
	}

	index := index.New()
	engine := engine.New(store, index)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repl := terminal.New(engine, cfg, config.Version)
	repl.Start()
}
