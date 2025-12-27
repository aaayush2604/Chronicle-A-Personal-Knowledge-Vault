package store

import (
	"bufio"
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

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		e, ok := parseLine(line)
		if !ok {
			continue
		}

		s.entries = append(s.entries, e)
		if e.ID >= s.nextID {
			s.nextID = e.ID + 1
		}
	}

	return scanner.Err()
}

func parseLine(line string) (entry.KnowledgeEntry, bool) {
	parts := strings.Split(line, "|")
	if len(parts) != 5 {
		return entry.KnowledgeEntry{}, false
	}

	version, _ := strconv.Atoi(parts[0])
	id, _ := strconv.Atoi(parts[1])
	ts, err := time.Parse(entry.TimeFormat, parts[2])
	if err != nil {
		return entry.KnowledgeEntry{}, false
	}

	return entry.KnowledgeEntry{
		Version:   entry.SchemaVersion(version),
		ID:        id,
		Timestamp: ts,
		Type:      entry.EntryType(parts[3]),
		Content:   parts[4],
	}, true
}
