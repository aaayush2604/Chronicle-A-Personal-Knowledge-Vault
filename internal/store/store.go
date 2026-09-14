package store

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
	"slices"
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

func (s *Store) Add(content string, tags []string, t entry.EntryType) (entry.KnowledgeEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var Tags []string
	for _, tag := range tags {
		if !slices.Contains(Tags, tag) {
			Tags = append(Tags, tag)
		}
	}
	e := entry.New(s.nextID, content, Tags)
	e.Type = t
	s.nextID++

	if err := s.append(e); err != nil {
		return entry.KnowledgeEntry{}, err
	}

	s.entries = append(s.entries, e)
	return e, nil
}

func (s *Store) AddUpdate(original entry.KnowledgeEntry, tags []*lexer.Token, t entry.EntryType) (entry.KnowledgeEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updated := original
	updated.Version = entry.CurrentVersion
	updated.Timestamp = time.Now()

	merged := make([]string, 0, len(original.Tags)+len(tags))
	for _, tag := range original.Tags {
		if !slices.Contains(merged, tag) {
			merged = append(merged, tag)
		}
	}
	for _, l := range tags {
		tag := l.Literal.(string)
		if !slices.Contains(merged, tag) {
			merged = append(merged, tag)
		}
	}
	updated.Tags = merged

	if t != "" {
		updated.Type = t
	}

	if err := s.append(updated); err != nil {
		return entry.KnowledgeEntry{}, err
	}

	update(s.entries, updated)
	return updated, nil
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

func update(list []entry.KnowledgeEntry, e entry.KnowledgeEntry) {
	flag := false
	for i, l := range list {
		if l.ID == e.ID {
			list[i] = e
			flag = true
		}
	}
	if !flag {
		list = append(list, e)
	}
}
