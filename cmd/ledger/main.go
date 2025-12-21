package main

import (
	"chronicle/internal/entry"
	"fmt"
	"os"
)

var entries []entry.KnowledgeEntry
var nextID int = 1

func main() {
	args := os.Args

	if len(args) < 2 {
		printUsage()
		return
	}

	var command string = args[1]

	switch command {
	case "add":
		handleAdd(args)
	case "list":
		handleList()
	default:
		fmt.Println("Unknown command:", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: ledger add <text>")
	fmt.Println("       ledger list")
}

func handleAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ledger add <text>")
		return
	}

	var content string = args[2]

	e := entry.New(nextID, content)
	nextID++

	entries = append(entries, e)

	fmt.Println("Saved entry with ID:", e.ID)
}

func handleList() {
	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	for _, e := range entries {
		fmt.Println("ID:", e.ID)
		fmt.Println("Timestamp:", e.Timestamp)
		fmt.Println("Content:", e.Content)
		fmt.Println("-----")
	}
}
