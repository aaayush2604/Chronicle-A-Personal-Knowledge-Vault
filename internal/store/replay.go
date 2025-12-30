package store

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"chronicle/internal/entry"
)

func (s *Store) replay() error {
	file, err := os.Open(s.logPath)
	if err != nil {
		return nil
	}

	defer file.Close()

	var warnings int
	lineNo := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineNo++
		line := scanner.Text()

		parts := strings.Split(line, "|")

		// DELETE record
		if len(parts) == 4 && parts[3] == "DEL" {
			id, _ := strconv.Atoi(parts[1])
			ts, _ := time.Parse(entry.TimeFormat, parts[2])
			s.deleted[id] = DeletionInfo{Timestamp: ts}
			continue
		}

		e, err := parseLine(line)
		if err != nil {
			warnings++
			continue
		}

		s.entries = append(s.entries, e)
		if e.ID >= s.nextID {
			s.nextID = e.ID + 1
		}
	}

	fmt.Printf(
		"Loaded %d entries\n⚠ %d entries could not be read and were skipped\n\n",
		lineNo-warnings,
		warnings,
	)

	return scanner.Err()
}

func parseLine(line string) (entry.KnowledgeEntry, error) {
	parts := strings.Split(line, "|")

	if len(parts) != 5 {
		return entry.KnowledgeEntry{}, fmt.Errorf("invalid field count")
	}

	version, _ := strconv.Atoi(parts[0])
	id, _ := strconv.Atoi(parts[1])
	ts, err := time.Parse(entry.TimeFormat, parts[2])
	if err != nil {
		return entry.KnowledgeEntry{}, fmt.Errorf("invalid time format ")
	}

	return entry.KnowledgeEntry{
		Version:   entry.SchemaVersion(version),
		ID:        id,
		Timestamp: ts,
		Type:      entry.EntryType(parts[3]),
		Content:   parts[4],
	}, nil
}
