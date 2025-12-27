package engine

import (
	"chronicle/internal/entry"
	"chronicle/internal/index"
)

func (e *Engine) Recall(term string) []entry.KnowledgeEntry {
	ids := index.Rank(e.index.Search(term))
	if len(ids) == 0 {
		return nil
	}

	entries := e.store.List()
	result := make([]entry.KnowledgeEntry, 0)

	idSet := make(map[int]struct{})
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	for _, e := range entries {
		if _, ok := idSet[e.ID]; ok {
			result = append(result, e)
		}
	}

	return result
}
