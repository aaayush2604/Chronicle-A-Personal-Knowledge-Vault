package store

import (
	"chronicle/internal/entry"
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	return s
}

func TestAddSingleEntry(t *testing.T) {
	s := newTestStore(t)

	e, err := s.Add("hello world", entry.TypeNote)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if e.ID != 1 {
		t.Fatalf("expected ID 1, got %d", e.ID)
	}

	if e.Content != "hello world" {
		t.Fatalf("expected content hello world, got %s", e.Content)
	}

	if e.Type != entry.TypeNote {
		t.Fatalf("expected note type")
	}

	if len(s.List()) != 1 {
		t.Fatalf("expected 1 entry")
	}
}

func TestAddMultipleEntries(t *testing.T) {
	s := newTestStore(t)

	e1, _ := s.Add("a", entry.TypeNote)
	e2, _ := s.Add("b", entry.TypeIdea)
	e3, _ := s.Add("c", entry.TypeQuestion)

	if e1.ID != 1 {
		t.Fatalf("expected first ID 1")
	}

	if e2.ID != 2 {
		t.Fatalf("expected second ID 2")
	}

	if e3.ID != 3 {
		t.Fatalf("expected third ID 3")
	}
}

func TestDeleteEntry(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("delete me", entry.TypeNote)

	err := s.Delete(e.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected deleted entry not to appear")
	}
}

func TestDeleteNonExistentEntry(t *testing.T) {
	s := newTestStore(t)

	err := s.Delete(999)

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDeleteAlreadyDeletedEntry(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("x", entry.TypeNote)

	if err := s.Delete(e.ID); err != nil {
		t.Fatal(err)
	}

	err := s.Delete(e.ID)

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCheckDelete(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("x", entry.TypeNote)

	if err := s.Delete(e.ID); err != nil {
		t.Fatal(err)
	}

	deleted, _ := s.CheckDelete(e.ID)

	if !deleted {
		t.Fatalf("expected deleted=true")
	}
}

func TestReplaySingleEntry(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Add("hello", entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Content != "hello" {
		t.Fatalf("content mismatch")
	}
}

func TestReplayDeletedEntry(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	e, _ := s.Add("hello", entry.TypeNote)

	if err := s.Delete(e.ID); err != nil {
		t.Fatal(err)
	}

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(reloaded.List()) != 0 {
		t.Fatalf("deleted entry should not appear")
	}
}

func TestReplayComputesNextID(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, _ := New(logPath)

	s.Add("a", entry.TypeNote)
	s.Add("b", entry.TypeNote)
	s.Add("c", entry.TypeNote)

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	e, err := reloaded.Add("d", entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	if e.ID != 4 {
		t.Fatalf("expected next ID to be 4, got %d", e.ID)
	}
}

func TestReplayCorruptedLine(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	err := os.WriteFile(
		logPath,
		[]byte("this is garbage\n"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected no valid entries")
	}
}

func TestReplayPartiallyCorruptedLog(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	content := `2|1|2025-01-01T10:00:00Z|note|first
this is garbage
2|2|2025-01-01T11:00:00Z|idea|second
`

	err := os.WriteFile(logPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := s.List()

	if len(entries) != 2 {
		t.Fatalf(
			"expected 2 valid entries, got %d",
			len(entries),
		)
	}
}

// func TestContentContainingPipeCharacter(t *testing.T) {
// 	s := newTestStore(t)

// 	_, err := s.Add(
// 		"hello | world",
// 		entry.TypeNote,
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	reloaded, err := New(s.logPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	entries := reloaded.List()

// 	if len(entries) != 1 {
// 		t.Fatalf("expected one entry")
// 	}

// 	if entries[0].Content != "hello | world" {
// 		t.Fatalf(
// 			"expected content preserved, got '%s'",
// 			entries[0].Content,
// 		)
// 	}
// }
