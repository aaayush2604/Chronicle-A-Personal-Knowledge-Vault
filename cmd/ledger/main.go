package main

import (
	"fmt"
	"log"

	"chronicle/internal/config"
	"chronicle/internal/engine"
	"chronicle/internal/index"
	"chronicle/internal/store"
	"chronicle/internal/terminal"
	"github.com/joho/godotenv"
)

func main() {
	if err:=godotenv.Load(); err != nil{
		log.Println("No .env file found, using OS environment")
	}

	logPath, err := config.LogPath()
	fmt.Println("Log Stored in: ", logPath)
	if err != nil {
		log.Fatal(err)
	}

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
