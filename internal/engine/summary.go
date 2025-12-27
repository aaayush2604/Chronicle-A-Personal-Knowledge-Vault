package engine

import "chronicle/internal/entry"

func (e *Engine) SummaryByType(entries []entry.KnowledgeEntry) map[entry.EntryType]int {
	out := make(map[entry.EntryType]int)
	for _, e := range entries {
		out[e.Type]++
	}
	return out
}
