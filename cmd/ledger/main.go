package main

import (
	"bufio"
	"chronicle/internal/entry"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var entries []entry.KnowledgeEntry
var nextID int = 1
var keywordIndex map[string][]int
var mu sync.RWMutex

const logFilePath = "data/chronicle.log"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	replayLog()

	wg.Add(1)
	go startIndexRebuilder(ctx, &wg)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

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
	case "search":
		handleSearch(args)
	default:
		fmt.Println("Unknown command:", command)
		printUsage()
	}

	go func() {
		<-sigCh
		fmt.Println("Shutting Down.....")
		cancel()
	}()

	cancel()

	wg.Wait()
}

func startIndexRebuilder(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			rebuildIndex()
		}
	}
}

func rebuildIndex() {
	mu.Lock()
	defer mu.Unlock()

	newIndex := make(map[string][]int)

	for _, e := range entries {
		words := tokenize(e.Content)
		for _, w := range words {
			newIndex[w] = append(newIndex[w], e.ID)
		}
	}

	keywordIndex = newIndex
}

func printUsage() {
	fmt.Println("Usage: ledger add <text>")
	fmt.Println("       ledger list")
	fmt.Println(" 	    ledger search <keyword>")
}

func handleAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ledger add <text>")
		return
	}

	var content string = args[2]

	mu.Lock()
	defer mu.Unlock()

	e := entry.New(nextID, content)
	nextID++

	if err := appendToLog(e); err != nil {
		fmt.Println("Error saving entry to log:", err)
		return
	}

	entries = append(entries, e)

	indexEntry(e)

	fmt.Println("Saved entry with ID:", e.ID)
}

func handleList() {
	mu.RLock()
	defer mu.RUnlock()

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

func handleSearch(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ledger search <word>")
		return
	}

	word := strings.ToLower(args[2])

	mu.RLock()
	defer mu.RUnlock()

	ids, ok := keywordIndex[word]
	if !ok || len(ids) == 0 {
		fmt.Println("No entries found for:", word)
		return
	}

	for _, id := range ids {
		if e, ok := getEntryByID(id); ok {
			fmt.Println("ID:", e.ID)
			fmt.Println("Time:", e.Timestamp.Format(time.RFC3339))
			fmt.Println("Content:", e.Content)
			fmt.Println("-----")
		}
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

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func replayLog() {
	mu.Lock()
	defer mu.Unlock()

	entries = nil
	nextID = 1
	keywordIndex = make(map[string][]int)

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

		indexEntry(e)

		if e.ID >= nextID {
			nextID = e.ID + 1
		}
	}
}

func indexEntry(e entry.KnowledgeEntry) {
	words := tokenize(e.Content)
	for _, w := range words {
		keywordIndex[w] = append(keywordIndex[w], e.ID)
	}
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func getEntryByID(id int) (entry.KnowledgeEntry, bool) {
	for _, e := range entries {
		if e.ID == id {
			return e, true
		}
	}
	return entry.KnowledgeEntry{}, false
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
