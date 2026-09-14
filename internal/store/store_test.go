package store

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
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

	e, err := s.Add("hello world", nil, entry.TypeNote)
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

	e1, _ := s.Add("a", nil, entry.TypeNote)
	e2, _ := s.Add("b", nil, entry.TypeIdea)
	e3, _ := s.Add("c", nil, entry.TypeQuestion)

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

	e, _ := s.Add("delete me", nil, entry.TypeNote)

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

	e, _ := s.Add("x", nil, entry.TypeNote)

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

	e, _ := s.Add("x", nil, entry.TypeNote)

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

	_, err = s.Add("hello", nil, entry.TypeNote)
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

	e, _ := s.Add("hello", nil, entry.TypeNote)

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

	s.Add("a", nil, entry.TypeNote)
	s.Add("b", nil, entry.TypeNote)
	s.Add("c", nil, entry.TypeNote)

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	e, err := reloaded.Add("d", nil, entry.TypeNote)
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

	content := `3|1||2025-01-01T10:00:00Z|note|first
this is garbage
3|2||2025-01-01T11:00:00Z|idea|second
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

func TestReplayPreservesTags(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	tags := []string{"go", "database"}

	_, err = s.Add(
		"hello",
		tags,
		entry.TypeNote,
	)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 1 {
		t.Fatalf("expected one entry")
	}

	if len(entries[0].Tags) != 2 {
		t.Fatalf("expected two tags")
	}

	if entries[0].Tags[0] != "go" {
		t.Fatalf("expected first tag go")
	}

	if entries[0].Tags[1] != "database" {
		t.Fatalf("expected second tag database")
	}
}

func TestReplayEntryWithoutTags(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Add(
		"hello",
		nil,
		entry.TypeNote,
	)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 1 {
		t.Fatalf("expected one entry")
	}

	if len(entries[0].Tags) != 0 {
		t.Fatalf("expected no tags")
	}
}

func TestSeparatorInContentSurvivesReplay(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	content := "pipe | inside | content"

	if _, err := s.Add(content, nil, entry.TypeNote); err != nil {
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

	if entries[0].Content != content {
		t.Fatalf("expected content %q, got %q", content, entries[0].Content)
	}
}

func TestVersion2EntryStillReplays(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	line := "2|1|2024-03-05T09:15:00+05:30|idea|an old v2 entry\n"
	if err := os.WriteFile(logPath, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := s.List()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Content != "an old v2 entry" {
		t.Fatalf("unexpected content %q", entries[0].Content)
	}

	if entries[0].Type != entry.TypeIdea {
		t.Fatalf("unexpected type %q", entries[0].Type)
	}
}

func TestAddUpdateKeepsTypeWhenNotSpecified(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("an idea", nil, entry.TypeIdea)

	updated, err := s.AddUpdate(e, []*lexer.Token{
		lexer.NewToken(lexer.TAG, "#go", "go", 0),
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	if updated.Type != entry.TypeIdea {
		t.Fatalf("expected type to carry over as idea, got %q", updated.Type)
	}
}

func TestAddUpdateOverridesTypeWhenSpecified(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("an idea", nil, entry.TypeIdea)

	updated, err := s.AddUpdate(e, nil, entry.TypeQuestion)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Type != entry.TypeQuestion {
		t.Fatalf("expected type question, got %q", updated.Type)
	}
}

func TestAddUpdateKeepsContentAndID(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("content that must survive", nil, entry.TypeLearning)

	updated, err := s.AddUpdate(e, nil, entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	if updated.ID != e.ID {
		t.Fatalf("expected id %d, got %d", e.ID, updated.ID)
	}

	if updated.Content != e.Content {
		t.Fatalf("expected content %q, got %q", e.Content, updated.Content)
	}
}

func TestAddUpdateDoesNotDuplicateTags(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("tagged", []string{"go", "db"}, entry.TypeNote)

	updated, err := s.AddUpdate(e, []*lexer.Token{
		lexer.NewToken(lexer.TAG, "#go", "go", 0),
		lexer.NewToken(lexer.TAG, "#ml", "ml", 0),
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"go", "db", "ml"}

	if len(updated.Tags) != len(expected) {
		t.Fatalf("expected tags %v, got %v", expected, updated.Tags)
	}

	for i, tag := range expected {
		if updated.Tags[i] != tag {
			t.Fatalf("expected tags %v, got %v", expected, updated.Tags)
		}
	}
}

func TestAddUpdateDoesNotMutateOriginalTags(t *testing.T) {
	s := newTestStore(t)

	e, _ := s.Add("tagged", []string{"a", "b", "c"}, entry.TypeNote)

	snapshot := s.List()[0]
	backing := snapshot.Tags[:cap(snapshot.Tags)]
	before := append([]string{}, backing...)

	if _, err := s.AddUpdate(snapshot, []*lexer.Token{
		lexer.NewToken(lexer.TAG, "#d", "d", 0),
	}, ""); err != nil {
		t.Fatal(err)
	}

	for i := range before {
		if backing[i] != before[i] {
			t.Fatalf("AddUpdate wrote into the caller's backing array: %v became %v", before, backing)
		}
	}

	_ = e
}

func TestUnreadableLogIsAnErrorNotAnEmptyVault(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	if err := os.WriteFile(logPath, []byte("3|1||2024-03-05T09:15:00+05:30|note|hidden\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(logPath, 0222); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(logPath, 0644)

	if _, err := New(logPath); err == nil {
		t.Fatalf("expected an unreadable log to be reported rather than read as an empty vault")
	}
}

func TestMissingLogIsNotAnError(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatalf("a log that does not exist yet should not be an error, got %v", err)
	}

	if len(s.List()) != 0 {
		t.Fatalf("expected an empty store")
	}
}

func TestReplayCollapsesRevisionsToLatest(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	lines := "" +
		"3|1||2024-03-01T09:00:00+05:30|note|first\n" +
		"3|2||2024-03-02T09:00:00+05:30|note|second\n" +
		"3|1|go|2024-03-03T09:00:00+05:30|idea|first revised\n" +
		"3|3||2024-03-04T09:00:00+05:30|note|written after the revision\n" +
		"3|4||2024-03-05T09:00:00+05:30|note|and another\n"

	if err := os.WriteFile(logPath, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := s.List()

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries after collapsing the revision, got %d", len(entries))
	}

	seen := map[int]entry.KnowledgeEntry{}
	for _, e := range entries {
		if _, duplicate := seen[e.ID]; duplicate {
			t.Fatalf("id %d appeared twice, the revision was not collapsed", e.ID)
		}
		seen[e.ID] = e
	}

	revised := seen[1]
	if revised.Content != "first revised" || revised.Type != entry.TypeIdea {
		t.Fatalf("expected id 1 to hold its latest version, got %q / %q", revised.Content, revised.Type)
	}

	if len(revised.Tags) != 1 || revised.Tags[0] != "go" {
		t.Fatalf("expected id 1 to hold its latest tags, got %v", revised.Tags)
	}

	for _, id := range []int{3, 4} {
		if _, ok := seen[id]; !ok {
			t.Fatalf("entry %d, written after the revision, was dropped", id)
		}
	}
}

func TestReplayAdvancesNextIDPastEveryRow(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	lines := "" +
		"3|1||2024-03-01T09:00:00+05:30|note|first\n" +
		"3|2||2024-03-02T09:00:00+05:30|note|second\n" +
		"3|1|go|2024-03-03T09:00:00+05:30|idea|first revised\n" +
		"3|3||2024-03-04T09:00:00+05:30|note|written after the revision\n" +
		"3|4||2024-03-05T09:00:00+05:30|note|and another\n"

	if err := os.WriteFile(logPath, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	e, err := s.Add("brand new", nil, entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	if e.ID != 5 {
		t.Fatalf("expected the next id to be 5, got %d. A reused id silently merges two entries on the next restart", e.ID)
	}
}

func TestRevisionWrittenLiveCollapsesOnRestart(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	s, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	first, _ := s.Add("keeps its text", nil, entry.TypeIdea)
	s.Add("a second entry", nil, entry.TypeNote)

	if _, err := s.AddUpdate(first, []*lexer.Token{
		lexer.NewToken(lexer.TAG, "#go", "go", 0),
	}, ""); err != nil {
		t.Fatal(err)
	}

	s.Add("written after the revision", nil, entry.TypeNote)

	reloaded, err := New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	entries := reloaded.List()

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries after restart, got %d", len(entries))
	}

	for _, e := range entries {
		if e.ID == first.ID {
			if e.Type != entry.TypeIdea {
				t.Fatalf("expected the revised entry to keep its type, got %q", e.Type)
			}
			if len(e.Tags) != 1 || e.Tags[0] != "go" {
				t.Fatalf("expected the revised entry to carry its new tag, got %v", e.Tags)
			}
		}
	}

	next, err := reloaded.Add("one more", nil, entry.TypeNote)
	if err != nil {
		t.Fatal(err)
	}

	if next.ID != 4 {
		t.Fatalf("expected the next id to be 4, got %d", next.ID)
	}
}
