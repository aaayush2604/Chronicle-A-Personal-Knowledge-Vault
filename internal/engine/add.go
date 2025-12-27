package engine

import "chronicle/internal/entry"

func (e *Engine) AddNote(content string, t entry.EntryType) (entry.KnowledgeEntry, error) {
	ke, err := e.store.Add(content, t)
	if err != nil {
		return entry.KnowledgeEntry{}, err
	}
	e.index.Build(e.store.List())

	return ke, nil
}
