package store

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"chronicle/internal/entry"
	"chronicle/internal/errorC"
)

func (s *Store) replay() error {
	file, err := os.Open(s.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errorC.New(errorC.Execution, fmt.Sprintf("Cannot read the log file at %s: %v", s.logPath, err))
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
			fmt.Printf("%s\n", err.Error())
			warnings++
			continue
		}

		update := false
		for i, c := range s.entries {
			if c.ID == e.ID {
				s.entries[i] = e
				update = true
				break
			}
		}

		if update {
			continue
		}

		s.entries = append(s.entries, e)
		if e.ID >= s.nextID {
			s.nextID = e.ID + 1
		}
	}

	fmt.Printf(
		"Read %d records from the log, %d entries loaded, %d deleted\n⚠ %d records could not be read and were skipped\n\n",
		lineNo,
		len(s.entries)-len(s.deleted),
		len(s.deleted),
		warnings,
	)

	return scanner.Err()
}

func parseLine(line string) (entry.KnowledgeEntry, error) {
	schema, _, _ := strings.Cut(line, "|")

	var version int
	var id int
	var tags []string
	var ts time.Time
	var err error

	if schema == "3" {
		parts := strings.SplitN(line, "|", 6)
		if len(parts) != 6 {
			fmt.Println("^")
			return entry.KnowledgeEntry{}, fmt.Errorf("invalid field count")
		}
		version, _ = strconv.Atoi(parts[0])
		id, _ = strconv.Atoi(parts[1])
		if parts[2] != "" {
			tags = strings.Split(parts[2], ",")
		}
		ts, err = time.Parse(entry.TimeFormat, parts[3])
		if err != nil {
			return entry.KnowledgeEntry{}, fmt.Errorf("invalid time format ")
		}

		return entry.KnowledgeEntry{
			Version:   entry.SchemaVersion(version),
			ID:        id,
			Tags:      tags,
			Timestamp: ts,
			Type:      entry.Canonical(parts[4]),
			Content:   parts[5],
		}, nil
	}

	if schema == "2" {
		parts := strings.SplitN(line, "|", 5)
		if len(parts) != 5 {
			fmt.Println("^")
			return entry.KnowledgeEntry{}, fmt.Errorf("invalid field count")
		}
		version, _ = strconv.Atoi(parts[0])
		id, _ = strconv.Atoi(parts[1])
		ts, err = time.Parse(entry.TimeFormat, parts[2])
		if err != nil {
			return entry.KnowledgeEntry{}, fmt.Errorf("invalid time format ")
		}

		return entry.KnowledgeEntry{
			Version:   entry.SchemaVersion(version),
			ID:        id,
			Tags:      tags,
			Timestamp: ts,
			Type:      entry.Canonical(parts[3]),
			Content:   parts[4],
		}, nil
	}

	return entry.KnowledgeEntry{}, errorC.New(errorC.Execution, fmt.Sprintf("Error in Parsing Entry from version: %s", schema))
}
