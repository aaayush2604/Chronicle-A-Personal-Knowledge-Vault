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

func TestEngineProcessDeletionDeleteAll(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	_, _ = eng.AddNote("one", entry.TypeNote)
	_, _ = eng.AddNote("two", entry.TypeNote)
	_, _ = eng.AddNote("three", entry.TypeNote)

	entries := s.List()

	deleted, err := eng.ProcessDeletion(nil, entries)
	if err != nil {
		t.Fatal(err)
	}

	if deleted != len(entries) {
		t.Fatalf("expected %d deletions, got %d", len(entries), deleted)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected all entries to be deleted")
	}
}

func TestEngineProcessDeletionSubset(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	e1, _ := eng.AddNote("one", entry.TypeNote)
	e2, _ := eng.AddNote("two", entry.TypeNote)
	e3, _ := eng.AddNote("three", entry.TypeNote)

	deleted, err := eng.ProcessDeletion([]int{e2.ID}, s.List())
	if err != nil {
		t.Fatal(err)
	}

	if deleted != 1 {
		t.Fatalf("expected 1 deletion, got %d", deleted)
	}

	remaining := s.List()

	if len(remaining) != 2 {
		t.Fatalf("expected 2 entries remaining")
	}

	for _, e := range remaining {
		if e.ID == e2.ID {
			t.Fatalf("entry should have been deleted")
		}
	}

	_ = e1
	_ = e3
}

func TestEngineProcessDeletionInvalidID(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	_, _ = eng.AddNote("one", entry.TypeNote)

	deleted, err := eng.ProcessDeletion([]int{999}, s.List())

	if err == nil {
		t.Fatalf("expected error")
	}

	if deleted != 0 {
		t.Fatalf("expected 0 deletions, got %d", deleted)
	}
}

func TestEngineProcessDeletionMixedIDs(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	e1, _ := eng.AddNote("one", entry.TypeNote)
	_, _ = eng.AddNote("two", entry.TypeNote)

	deleted, err := eng.ProcessDeletion([]int{e1.ID, 999}, s.List())

	if err == nil {
		t.Fatalf("expected error")
	}

	if deleted != 0 {
		t.Fatalf("expected 0 deletions when validation fails, got %d", deleted)
	}
}

func TestEngineProcessDeletionEmptyResult(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	deleted, err := eng.ProcessDeletion(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if deleted != 0 {
		t.Fatalf("expected 0 deletions, got %d", deleted)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected store to remain empty")
	}
}

func TestEngineForgetWorkflow(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	eng := engine.New(s, idx)

	_, _ = eng.Query(`rem "database systems"`)
	_, _ = eng.Query(`rem "machine learning"`)

	results, err := eng.Query(`forget contains["database"]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected one matching entry")
	}

	deleted, err := eng.ProcessDeletion(nil, results)
	if err != nil {
		t.Fatal(err)
	}

	if deleted != len(results) {
		t.Fatalf("expected %d deletions, got %d", len(results), deleted)
	}

	remaining, err := eng.Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(remaining) != 1 {
		t.Fatalf("expected one remaining entry")
	}

	if remaining[0].Content != "machine learning" {
		t.Fatalf("wrong entry was deleted")
	}
}
