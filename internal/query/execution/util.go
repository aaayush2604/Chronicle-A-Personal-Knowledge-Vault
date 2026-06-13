package execution

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"

	"strings"
)

func EntryToRecord(e entry.KnowledgeEntry) Record {
	return Record{
		"id":      e.ID,
		"content": e.Content,
		"date":    e.Timestamp,
		"len":     len(e.Content),
		"type":    string(e.Type),
	}
}

func TokensToString(list []*lexer.Token) string {
	var sb strings.Builder
	for _, l := range list {
		sb.WriteString(l.Lexeme)
		sb.WriteString(" ")
	}
	return sb.String()
}
