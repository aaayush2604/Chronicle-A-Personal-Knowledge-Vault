package store

import (
	"chronicle/internal/entry"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	entries []entry.KnowledgeEntry
	nextID  int
	logPath string
}

func New(logPath string) (*Store, error) {
	s := &Store{
		logPath: logPath,
		nextID:  1,
	}

	if err := s.replay(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Add(content string) (entry.KnowledgeEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := entry.New(s.nextID, content)
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

	out := make([]entry.KnowledgeEntry, len(s.entries))
	copy(out, s.entries)
	return out
}
