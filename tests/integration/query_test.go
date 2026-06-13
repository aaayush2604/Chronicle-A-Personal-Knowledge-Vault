package integration

import (
	"chronicle/internal/engine"
	"chronicle/internal/entry"
	"chronicle/internal/index"
	"chronicle/internal/store"
	"path/filepath"
	"testing"
)

func buildEngine(t *testing.T) *engine.Engine {
	t.Helper()

	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	idx := index.New()

	return engine.New(s, idx)
}

func TestQueryAll(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.AddNote(
		"note one",
		entry.TypeNote,
	)

	_, _ = eng.AddNote(
		"note two",
		entry.TypeIdea,
	)

	results, err := eng.Query(
		`recall all`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}
}

func TestQueryContains(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.AddNote(
		"golang database systems",
		entry.TypeNote,
	)

	_, _ = eng.AddNote(
		"machine learning",
		entry.TypeIdea,
	)

	results, err := eng.Query(
		`recall contains["database"]`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}
}

func TestQueryTypeFilter(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.AddNote(
		"note",
		entry.TypeNote,
	)

	_, _ = eng.AddNote(
		"idea",
		entry.TypeIdea,
	)

	results, err := eng.Query(
		`recall type["note"]`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	if results[0].Type != entry.TypeNote {
		t.Fatalf("expected note")
	}
}

func TestQueryComparison(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.AddNote(
		"this is a very long entry",
		entry.TypeNote,
	)

	results, err := eng.Query(
		`recall len > 10`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}
}

func TestQueryLogicalAnd(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.AddNote(
		"database systems",
		entry.TypeNote,
	)

	_, _ = eng.AddNote(
		"database systems",
		entry.TypeIdea,
	)

	results, err := eng.Query(
		`recall contains["database"] and type["note"]`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}
}

func TestInvalidQuery(t *testing.T) {
	eng := buildEngine(t)

	_, err := eng.Query(
		`len > 10`,
	)

	if err == nil {
		t.Fatalf("expected query error")
	}
}

func TestDeletedEntriesNotReturnedAfterRestart(t *testing.T) {
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

	reloadedStore, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	reloadedEngine := engine.New(
		reloadedStore,
		index.New(),
	)

	results, err := reloadedEngine.Query(
		`recall all`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatalf(
			"expected 0 results, got %d",
			len(results),
		)
	}
}

func TestQueryNoteCommand(t *testing.T) {
	eng := buildEngine(t)

	results, err := eng.Query(
		`note #go #database "hello world"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	if results[0].Content != `"hello world" ` {
		t.Fatalf("content mismatch")
	}

	if len(results[0].Tags) != 2 {
		t.Fatalf("expected two tags")
	}
}
