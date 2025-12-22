package main

import (
	"bufio"
	"chronicle/internal/entry"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var entries []entry.KnowledgeEntry
var nextID int = 1

const logFilePath = "data/chronicle.log"

func main() {
	loadFromLog()

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

	if err := appendToLog(e); err != nil {
		fmt.Println("Error saving entry to log:", err)
		return
	}

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

func appendToLog(e entry.KnowledgeEntry) error {
	file, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	line := fmt.Sprintf(
		"%d|%s|%s\n",
		e.ID,
		e.Timestamp.Format(time.RFC3339),
		e.Content,
	)

	if _, err := writer.WriteString(line); err != nil {
		return err
	}

	return writer.Flush()
}

func loadFromLog() {
	file, err := os.Open(logFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		e, ok := parseLogLine(line)
		if !ok {
			continue
		}
		entries = append(entries, e)

		if e.ID >= nextID {
			nextID = e.ID + 1
		}
	}
}

func parseLogLine(line string) (entry.KnowledgeEntry, bool) {
	parts := strings.Split(line, "|")
	if len(parts) != 3 {
		return entry.KnowledgeEntry{}, false
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return entry.KnowledgeEntry{}, false
	}
	timestamp, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return entry.KnowledgeEntry{}, false
	}
	content := parts[2]

	return entry.KnowledgeEntry{
		ID:        id,
		Timestamp: timestamp,
		Content:   content,
	}, true
}
