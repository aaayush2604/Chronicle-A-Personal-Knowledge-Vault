package store

import (
	"bufio"
	"chronicle/internal/entry"
	"fmt"
	"os"
	"strings"
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

	var sb strings.Builder
	for _, t := range e.Tags {
		sb.WriteString(t)
		sb.WriteString(",")
	}
	tags := sb.String()
	if len(tags) > 0 && tags[len(tags)-1] == ',' {
		tags = tags[:len(tags)-1]
	}
	writer := bufio.NewWriter(file)

	line := fmt.Sprintf(
		"%d|%d|%s|%s|%s|%s\n",
		e.Version,
		e.ID,
		tags,
		e.Timestamp.Format(entry.TimeFormat),
		e.Type,
		e.Content,
	)

	if _, err := writer.WriteString(line); err != nil {
		return err
	}

	return writer.Flush()
}
