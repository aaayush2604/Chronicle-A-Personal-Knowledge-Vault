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

func TestTypeForAcceptsNamesAndAliases(t *testing.T) {
	cases := map[string]EntryType{
		"note": TypeNote, "n": TypeNote,
		"idea": TypeIdea, "i": TypeIdea,
		"question": TypeQuestion, "q": TypeQuestion,
		"learning": TypeLearning, "l": TypeLearning,
		"important": TypeImportant, "imp": TypeImportant,
		"IMPORTANT": TypeImportant, "Imp": TypeImportant,
	}

	for name, want := range cases {
		got, ok := TypeFor(name)
		if !ok {
			t.Fatalf("%q: expected a known type", name)
		}
		if got != want {
			t.Fatalf("%q: expected %q, got %q", name, want, got)
		}
	}

	if _, ok := TypeFor("todo"); ok {
		t.Fatalf("expected todo to be unknown")
	}
}

func TestCanonicalUpgradesStoredAliases(t *testing.T) {
	if got := Canonical("imp"); got != TypeImportant {
		t.Fatalf("expected a stored imp to read back as important, got %q", got)
	}

	if got := Canonical("todo"); got != EntryType("todo") {
		t.Fatalf("expected an unknown stored type to be left alone, got %q", got)
	}
}
