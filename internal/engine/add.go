package engine

import "chronicle/internal/entry"

func (e *Engine) AddNote(content string, t entry.EntryType) (entry.KnowledgeEntry, error) {
	ke, err := e.store.Add(content)
	if err != nil {
		return entry.KnowledgeEntry{}, err
	}

	ke.Type = t

	e.index.Build(e.store.List())

	return ke, nil
}
