package entry

import (
	"testing"
	"time"
)

func TestNewEntrySetsID(t *testing.T) {
	e := New(42, "hello", []string{})

	if e.ID != 42 {
		t.Fatalf("expected ID 42, got %d", e.ID)
	}
}

func TestNewEntrySetsContent(t *testing.T) {
	e := New(1, "my content", []string{})

	if e.Content != "my content" {
		t.Fatalf("expected content 'my content', got '%s'", e.Content)
	}
}

func TestNewEntrySetsCurrentVersion(t *testing.T) {
	e := New(1, "hello", []string{})

	if e.Version != CurrentVersion {
		t.Fatalf(
			"expected version %v, got %v",
			CurrentVersion,
			e.Version,
		)
	}
}

func TestNewEntrySetsDefaultType(t *testing.T) {
	e := New(1, "hello", []string{})

	if e.Type != TypeNote {
		t.Fatalf(
			"expected type %v, got %v",
			TypeNote,
			e.Type,
		)
	}
}

func TestNewEntrySetsTimestamp(t *testing.T) {
	before := time.Now()

	e := New(1, "hello", []string{})

	after := time.Now()

	if e.Timestamp.Before(before) {
		t.Fatalf(
			"timestamp %v occurred before constructor call",
			e.Timestamp,
		)
	}

	if e.Timestamp.After(after) {
		t.Fatalf(
			"timestamp %v occurred after constructor completed",
			e.Timestamp,
		)
	}
}

func TestEntryTypesAreUnique(t *testing.T) {
	types := []EntryType{
		TypeNote,
		TypeIdea,
		TypeQuestion,
		TypeLearning,
		TypeImportant,
	}

	seen := make(map[EntryType]bool)

	for _, typ := range types {
		if seen[typ] {
			t.Fatalf("duplicate entry type found: %v", typ)
		}
		seen[typ] = true
	}
}

func TestCurrentVersionMatchesLatestVersion(t *testing.T) {
	if CurrentVersion != Version3 {
		t.Fatalf(
			"CurrentVersion should point to latest schema version",
		)
	}
}
