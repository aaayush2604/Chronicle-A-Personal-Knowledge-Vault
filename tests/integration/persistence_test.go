package integration

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
	"chronicle/internal/store"
	"path/filepath"
	"testing"
)

func TestPersistenceAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Add("first note", []*lexer.Token{}, entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Add("second note", []*lexer.Token{}, entry.TypeIdea)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Content != "first note" {
		t.Fatalf("first entry mismatch")
	}

	if entries[1].Content != "second note" {
		t.Fatalf("second entry mismatch")
	}
}

func TestDeletionPersistsAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	e, err := s.Add("delete me", []*lexer.Token{}, entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Delete(e.ID)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 0 {
		t.Fatalf("expected deleted entry to remain deleted")
	}
}
