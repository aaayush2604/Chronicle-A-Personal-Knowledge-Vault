package entry

import "time"

const TimeFormat = time.RFC3339

type KnowledgeEntry struct {
	Version   SchemaVersion
	ID        int
	Timestamp time.Time
	Type      EntryType
	Content   string
}

func New(id int, content string) KnowledgeEntry {
	return KnowledgeEntry{
		Version:   CurrentVersion,
		ID:        id,
		Timestamp: time.Now(),
		Type:      TypeNote,
		Content:   content,
	}
}
