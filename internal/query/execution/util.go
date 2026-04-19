package execution

import "chronicle/internal/entry"

func EntryToRecord(e entry.KnowledgeEntry) Record {
	return Record{
		"id":      e.ID,
		"content": e.Content,
		"date":    e.Timestamp,
		"len":     len(e.Content),
		"type":    string(e.Type),
	}
}
