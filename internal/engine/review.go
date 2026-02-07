package engine

import (
	"chronicle/internal/entry"
	"time"
)

func (e *Engine) Today() []entry.KnowledgeEntry {
	now := time.Now()
	return filterByDate(e.store.List(), startOfDay(now))
}

func (e *Engine) ThisWeek() []entry.KnowledgeEntry {
	now := time.Now()
	start := now.AddDate(0, 0, -7)
	return filterByDate(e.store.List(), start)
}

func (e *Engine) ThisMonth() []entry.KnowledgeEntry {
	now := time.Now()
	start := now.AddDate(0, -1, 0)
	return filterByDate(e.store.List(), start)
}

func (e *Engine) ThisYear() []entry.KnowledgeEntry {
	now := time.Now()
	start := now.AddDate(-1, 0, 0)
	return filterByDate(e.store.List(), start)
}

func filterByDate(entries []entry.KnowledgeEntry, from time.Time) []entry.KnowledgeEntry {
	out := make([]entry.KnowledgeEntry, 0)
	for _, e := range entries {
		if e.Timestamp.After(from) {
			out = append(out, e)
		}
	}
	return out
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
