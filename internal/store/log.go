package store

import (
	"bufio"
	"chronicle/internal/entry"
	"fmt"
	"os"
)

func (s *Store) append(e entry.KnowledgeEntry) error {
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

	line := fmt.Sprintf(
		"%d|%d|%s|%s|%s\n",
		e.Version,
		e.ID,
		e.Timestamp.Format(entry.TimeFormat),
		e.Type,
		e.Content,
	)

	if _, err := writer.WriteString(line); err != nil {
		return err
	}

	return writer.Flush()
}
