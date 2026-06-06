package integration

import (
	"chronicle/internal/engine"
	"chronicle/internal/entry"
	"chronicle/internal/index"
	"chronicle/internal/store"
	"path/filepath"
	"testing"
)

func TestEngineAddNote(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	e, err := eng.AddNote(
		"hello world",
		entry.TypeNote,
	)
	if err != nil {
		t.Fatal(err)
	}

	if e.ID != 1 {
		t.Fatalf("expected ID 1, got %d", e.ID)
	}

	entries := s.List()

	if len(entries) != 1 {
		t.Fatalf("expected one entry")
	}
}

func TestEngineDelete(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	e, err := eng.AddNote(
		"delete me",
		entry.TypeNote,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = eng.Delete(e.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected deleted entry not to appear")
	}
}

func TestEngineRecall(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	_, _ = eng.AddNote(
		"golang database systems",
		entry.TypeNote,
	)

	results := eng.Recall("database")

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 recall result, got %d",
			len(results),
		)
	}
}

func TestSummaryByType(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	entries := []entry.KnowledgeEntry{
		{Type: entry.TypeNote},
		{Type: entry.TypeNote},
		{Type: entry.TypeIdea},
	}

	summary := eng.SummaryByType(entries)

	if summary[entry.TypeNote] != 2 {
		t.Fatalf("expected 2 notes")
	}

	if summary[entry.TypeIdea] != 1 {
		t.Fatalf("expected 1 idea")
	}
}
