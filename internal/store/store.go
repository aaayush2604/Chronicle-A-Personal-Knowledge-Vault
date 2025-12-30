package store

import (
	"chronicle/internal/entry"
	"sync"
	"time"
)

type DeletionInfo struct {
	Timestamp time.Time
}

type Store struct {
	mu      sync.RWMutex
	entries []entry.KnowledgeEntry
	deleted map[int]DeletionInfo
	nextID  int
	logPath string
}

func New(logPath string) (*Store, error) {
	s := &Store{
		logPath: logPath,
		nextID:  1,
		deleted: make(map[int]DeletionInfo),
	}

	if err := s.replay(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Add(content string, t entry.EntryType) (entry.KnowledgeEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := entry.New(s.nextID, content)
	e.Type = t
	s.nextID++

	if err := s.append(e); err != nil {
		return entry.KnowledgeEntry{}, err
	}

	s.entries = append(s.entries, e)
	return e, nil
}

func (s *Store) List() []entry.KnowledgeEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []entry.KnowledgeEntry
	for _, e := range s.entries {
		if _, deleted := s.deleted[e.ID]; deleted {
			continue
		}
		out = append(out, e)
	}
	return out
}
