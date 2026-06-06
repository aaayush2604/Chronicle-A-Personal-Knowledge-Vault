package index

import (
	"chronicle/internal/entry"
	"testing"
)

func TestEmptyIndexSearch(t *testing.T) {
	idx := New()

	results := idx.Search("golang")

	if len(results) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestBuildSingleEntry(t *testing.T) {
	idx := New()

	entries := []entry.KnowledgeEntry{
		{
			ID:      1,
			Content: "golang database",
		},
	}

	idx.Build(entries)

	results := idx.Search("golang")

	if len(results) != 1 {
		t.Fatalf("expected one result")
	}

	if results[0] != 1 {
		t.Fatalf("expected id 1")
	}
}

func TestBuildMultipleEntries(t *testing.T) {
	idx := New()

	entries := []entry.KnowledgeEntry{
		{
			ID:      1,
			Content: "golang database",
		},
		{
			ID:      2,
			Content: "database systems",
		},
	}

	idx.Build(entries)

	results := idx.Search("database")

	if len(results) != 2 {
		t.Fatalf("expected two hits")
	}
}

func TestNormalizeCase(t *testing.T) {
	idx := New()

	entries := []entry.KnowledgeEntry{
		{
			ID:      1,
			Content: "GoLang",
		},
	}

	idx.Build(entries)

	results := idx.Search("golang")

	if len(results) != 1 {
		t.Fatalf("expected normalized search")
	}
}

func TestStopWordsRemoved(t *testing.T) {
	idx := New()

	entries := []entry.KnowledgeEntry{
		{
			ID:      1,
			Content: "the and is of",
		},
	}

	idx.Build(entries)

	results := idx.Search("the")

	if len(results) != 0 {
		t.Fatalf("stopwords should not be indexed")
	}
}

func TestReset(t *testing.T) {
	idx := New()

	entries := []entry.KnowledgeEntry{
		{
			ID:      1,
			Content: "golang",
		},
	}

	idx.Build(entries)

	idx.Reset()

	results := idx.Search("golang")

	if len(results) != 0 {
		t.Fatalf("expected empty index after reset")
	}
}

func TestRank(t *testing.T) {
	ids := []int{
		1,
		1,
		1,
		2,
		2,
		3,
	}

	ranked := Rank(ids)

	if ranked[0] != 1 {
		t.Fatalf("expected most frequent id first")
	}
}
