package entry

import "time"

type KnowledgeEntry struct {
	ID        int
	Timestamp time.Time
	Content   string
}

func New(id int, content string) KnowledgeEntry {
	return KnowledgeEntry{
		ID:        id,
		Timestamp: time.Now(),
		Content:   content,
	}
}
