package store

import (
	"bufio"
	"chronicle/internal/entry"
	"fmt"
	"os"
	"time"
)

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.deleted[id]; ok {
		return fmt.Errorf("entry already deleted")
	}

	found := false
	var version entry.SchemaVersion
	for _, e := range s.entries {
		if e.ID == id {
			found = true
			version = e.Version
			break
		}
	}
	if !found {
		return fmt.Errorf("Entry not found")
	}

	ts := time.Now()

	line := fmt.Sprintf(
		"%d|%d|%s|DEL\n",
		version,
		id,
		ts.Format(entry.TimeFormat),
	)

	file, err := os.OpenFile(
		s.logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString(line); err != nil {
		return err
	}

	s.deleted[id] = DeletionInfo{Timestamp: ts}
	return writer.Flush()
}

func (s *Store) CheckDelete(id int) (bool, DeletionInfo) {
	ts, ok := s.deleted[id]
	if ok {
		return true, ts
	}
	return false, DeletionInfo{}
}
