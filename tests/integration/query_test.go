package integration

import (
	"chronicle/internal/engine"
	"chronicle/internal/entry"
	"chronicle/internal/index"
	"chronicle/internal/store"
	"os"
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
		`recall type[note]`,
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
		`recall contains["database"] and type[note]`,
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

func TestQueryRemCommand(t *testing.T) {
	eng := buildEngine(t)

	results, err := eng.Query(
		`rem #go #database "hello world"`,
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

	if results[0].Content != "hello world" {
		t.Fatalf("content mismatch")
	}

	if len(results[0].Tags) != 2 {
		t.Fatalf("expected two tags")
	}
}

func TestQueryTags(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(
		`rem #go "golang project"`,
	)

	_, _ = eng.Query(
		`rem #python "python project"`,
	)

	results, err := eng.Query(
		`recall tags[go]`,
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

func TestQueryTagsOr(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(
		`rem #go "golang project"`,
	)

	_, _ = eng.Query(
		`rem #database "sql project"`,
	)

	results, err := eng.Query(
		`recall tags[go,database]`,
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

func TestRemDefaultType(t *testing.T) {
	eng := buildEngine(t)

	results, err := eng.Query(
		`rem "hello world"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].Type != entry.TypeNote {
		t.Fatalf("expected note type")
	}
}

func TestRemLearningType(t *testing.T) {
	eng := buildEngine(t)

	results, err := eng.Query(
		`rem @learning "hello world"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].Type != entry.TypeLearning {
		t.Fatalf("expected learning type")
	}
}

func TestRemTags(t *testing.T) {
	eng := buildEngine(t)

	results, err := eng.Query(
		`rem @learning #go #database "hello world"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results[0].Tags) != 2 {
		t.Fatalf("expected two tags")
	}

	if results[0].Type != entry.TypeLearning {
		t.Fatalf("expected learning type")
	}

	if results[0].Content != "hello world" {
		t.Fatalf("content mismatch")
	}
}

func TestQueryForgetAll(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem "first"`)
	_, _ = eng.Query(`rem "second"`)

	results, err := eng.Query(`forget all`)
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

func TestQueryForgetContains(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem "golang database"`)

	_, _ = eng.Query(`rem "machine learning"`)

	results, err := eng.Query(
		`forget contains["database"]`,
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

	if results[0].Content != "golang database" {
		t.Fatalf("unexpected entry returned")
	}
}

func TestQueryForgetType(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note "note"`)

	_, _ = eng.Query(`rem @idea "idea"`)

	results, err := eng.Query(
		`forget type[note]`,
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

func TestQueryForgetLogicalExpression(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(
		`rem @note "database systems"`,
	)

	_, _ = eng.Query(
		`rem @idea "database systems"`,
	)

	results, err := eng.Query(
		`forget contains["database"] and type[note]`,
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

func TestReviseKeepsTypeWhenOnlyTagsChange(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @idea #go an idea worth keeping`)

	results, err := eng.Query(`revise #db where tags[go]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Type != entry.TypeIdea {
		t.Fatalf("expected type idea to carry over, got %q", results[0].Type)
	}

	if results[0].Content != "an idea worth keeping" {
		t.Fatalf("unexpected content %q", results[0].Content)
	}
}

func TestReviseIsIdempotentForTags(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note #go a note`)
	_, _ = eng.Query(`revise #go where tags[go]`)

	results, err := eng.Query(`revise #go where tags[go]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results[0].Tags) != 1 {
		t.Fatalf("expected a single tag, got %v", results[0].Tags)
	}
}

func dateTimeEngine(t *testing.T) *engine.Engine {
	t.Helper()

	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	lines := "" +
		"3|1||2026-09-13T19:30:00+05:30|note|evening\n" +
		"3|2||2026-09-13T09:00:00+05:30|note|morning\n" +
		"3|3||2025-01-05T10:00:00+05:30|idea|last year\n"

	if err := os.WriteFile(logPath, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	return engine.New(s, index.New())
}

func TestDateTimeComparisons(t *testing.T) {
	eng := dateTimeEngine(t)

	cases := []struct {
		query string
		ids   []int
	}{
		{`recall time > "7 PM"`, []int{1}},
		{`recall time >= "7 PM"`, []int{1}},
		{`recall time = "7 PM"`, []int{1}},
		{`recall time < "7 PM"`, []int{2, 3}},
		{`recall time > "19:00"`, []int{1}},
		{`recall date = "2026-09-13"`, []int{1, 2}},
		{`recall date = "13/09/2026"`, []int{1, 2}},
		{`recall date > "2026-09-13"`, []int{}},
		{`recall date < "2026-09-13"`, []int{3}},
		{`recall date >= "2026-01-01"`, []int{1, 2}},
		{`recall date = "2026"`, []int{1, 2}},
		{`recall date = "2026-09"`, []int{1, 2}},
		{`recall len > 3.5`, []int{1, 2, 3}},
		{`recall len > 100.5`, []int{}},
	}

	for _, c := range cases {
		results, err := eng.Query(c.query)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", c.query, err)
		}

		if len(results) != len(c.ids) {
			t.Fatalf("%s: expected ids %v, got %d results", c.query, c.ids, len(results))
		}

		for i, id := range c.ids {
			if results[i].ID != id {
				t.Fatalf("%s: expected ids %v, got a result with id %d", c.query, c.ids, results[i].ID)
			}
		}
	}
}

func TestInvalidDateTimeValuesAreRejected(t *testing.T) {
	eng := dateTimeEngine(t)

	for _, query := range []string{
		`recall date > "yesterday"`,
		`recall date = "2026-13-45"`,
		`recall time > "25:00"`,
		`recall time = "13 PM"`,
	} {
		if _, err := eng.Query(query); err == nil {
			t.Fatalf("%s: expected an error instead of matching everything", query)
		}
	}
}

func TestImportantIsFilterableByTypeAndAlias(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "chronicle.log")

	legacy := "3|1||2024-03-05T09:15:00+05:30|imp|stored back when the type was imp\n"
	if err := os.WriteFile(logPath, []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	eng := engine.New(s, index.New())

	if _, err := eng.Query(`rem @imp written through the alias`); err != nil {
		t.Fatal(err)
	}

	if _, err := eng.Query(`rem @important written through the full name`); err != nil {
		t.Fatal(err)
	}

	for _, query := range []string{`recall type[imp]`, `recall type[important]`} {
		results, err := eng.Query(query)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", query, err)
		}

		if len(results) != 3 {
			t.Fatalf("%s: expected 3 results, got %d", query, len(results))
		}

		for _, r := range results {
			if r.Type != entry.TypeImportant {
				t.Fatalf("%s: expected every result to be important, got %q", query, r.Type)
			}
		}
	}
}

func TestContainsListIsOr(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note a machine learning project`)
	_, _ = eng.Query(`rem @idea a golang database`)

	cases := []struct {
		query string
		count int
	}{
		{`recall contains["learning","golang"]`, 2},
		{`recall contains["golang","zzz"]`, 1},
		{`recall contains["zzz","qqq"]`, 0},
		{`recall contains["project"] AND contains["learning"]`, 1},
		{`recall contains["project"] AND contains["golang"]`, 0},
	}

	for _, c := range cases {
		results, err := eng.Query(c.query)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", c.query, err)
		}

		if len(results) != c.count {
			t.Fatalf("%s: expected %d results, got %d", c.query, c.count, len(results))
		}
	}
}

func openEngine(t *testing.T, logPath string) *engine.Engine {
	t.Helper()

	s, err := store.New(logPath)
	if err != nil {
		t.Fatal(err)
	}

	return engine.New(s, index.New())
}

func TestRemRoundTripsAcrossRestart(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	if _, err := eng.Query(`rem #go @idea Meeting with John, at 5 PM. cost 50%`); err != nil {
		t.Fatal(err)
	}

	results, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(results))
	}

	got := results[0]

	if got.Content != "Meeting with John, at 5 PM. cost 50%" {
		t.Fatalf("content was not stored verbatim, got %q", got.Content)
	}

	if got.Type != entry.TypeIdea {
		t.Fatalf("expected idea, got %q", got.Type)
	}

	if len(got.Tags) != 1 || got.Tags[0] != "go" {
		t.Fatalf("expected tags [go], got %v", got.Tags)
	}
}

func TestRemRejectsBadInput(t *testing.T) {
	eng := buildEngine(t)

	for _, query := range []string{
		`rem @todo buy milk`,
		`rem @idea`,
		`rem @idea #go`,
	} {
		if _, err := eng.Query(query); err == nil {
			t.Fatalf("%s: expected an error", query)
		}
	}
}

func TestRecallWithGroupingEvaluates(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note a golang note`)
	_, _ = eng.Query(`rem @idea a golang idea`)
	_, _ = eng.Query(`rem @question a rust question`)

	cases := []struct {
		query string
		count int
	}{
		{`recall (type[note])`, 1},
		{`recall ((type[note]))`, 1},
		{`recall (type[note] OR type[idea])`, 2},
		{`recall (type[note] OR type[idea]) AND contains["golang"]`, 2},
		{`recall (type[note] OR type[question]) AND contains["rust"]`, 1},
		{`recall ((type[note] OR type[idea]) AND contains["golang"]) OR type[question]`, 3},
	}

	for _, c := range cases {
		results, err := eng.Query(c.query)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", c.query, err)
		}

		if len(results) != c.count {
			t.Fatalf("%s: expected %d results, got %d", c.query, c.count, len(results))
		}
	}
}

func TestForgetQueryAloneDeletesNothing(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @note second`)

	matched, err := eng.Query(`forget type[note]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(matched) != 2 {
		t.Fatalf("expected the forget query to match 2 entries, got %d", len(matched))
	}

	remaining, err := eng.Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(remaining) != 2 {
		t.Fatalf("a forget query on its own must not delete, still expected 2 entries, got %d", len(remaining))
	}
}

func TestForgetSubsetLeavesTheRest(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @note second`)
	_, _ = eng.Query(`rem @note third`)

	matched, err := eng.Query(`forget type[note]`)
	if err != nil {
		t.Fatal(err)
	}

	deleted, err := eng.ProcessDeletion([]int{matched[1].ID}, matched)
	if err != nil {
		t.Fatal(err)
	}

	if deleted != 1 {
		t.Fatalf("expected 1 deletion, got %d", deleted)
	}

	remaining, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(remaining) != 2 {
		t.Fatalf("expected 2 entries after restart, got %d", len(remaining))
	}

	for _, e := range remaining {
		if e.Content == "second" {
			t.Fatalf("the deleted entry came back after a restart")
		}
	}
}

func TestReviseSurvivesRestart(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	if _, err := eng.Query(`rem @note #draft something worth keeping`); err != nil {
		t.Fatal(err)
	}

	if _, err := eng.Query(`revise @important #final where tags[draft]`); err != nil {
		t.Fatal(err)
	}

	results, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 entry after restart, got %d", len(results))
	}

	got := results[0]

	if got.Type != entry.TypeImportant {
		t.Fatalf("expected the new type to persist, got %q", got.Type)
	}

	if got.Content != "something worth keeping" {
		t.Fatalf("revise must not touch the text, got %q", got.Content)
	}

	if len(got.Tags) != 2 || got.Tags[0] != "draft" || got.Tags[1] != "final" {
		t.Fatalf("expected tags [draft final], got %v", got.Tags)
	}
}

func TestReviseWhereAllUpdatesEveryEntry(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @idea second`)
	_, _ = eng.Query(`rem @question third`)

	results, err := eng.Query(`revise @learning where all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 entries revised, got %d", len(results))
	}

	after, err := eng.Query(`recall type[learning]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(after) != 3 {
		t.Fatalf("expected every entry to be a learning, got %d", len(after))
	}
}

func TestReviseMatchingNothingChangesNothing(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	_, _ = eng.Query(`rem @note untouched`)

	results, err := eng.Query(`revise @important where tags[missing]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatalf("expected nothing to be revised, got %d", len(results))
	}

	after, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(after) != 1 || after[0].Type != entry.TypeNote {
		t.Fatalf("expected the entry to be untouched, got %d entries of type %q", len(after), after[0].Type)
	}
}

func TestReviseUpdatesTheEntryTime(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	legacy := "3|1|old|2024-03-05T09:15:00+05:30|idea|written long ago\n"
	if err := os.WriteFile(logPath, []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}

	eng := openEngine(t, logPath)

	before, err := eng.Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	results, err := eng.Query(`revise #fresh where tags[old]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 entry revised, got %d", len(results))
	}

	if !results[0].Timestamp.After(before[0].Timestamp) {
		t.Fatalf("expected revise to move the time forward, was %v now %v", before[0].Timestamp, results[0].Timestamp)
	}

	if results[0].Type != entry.TypeIdea {
		t.Fatalf("expected the type to carry over, got %q", results[0].Type)
	}
}

func TestReviseDoesNotResurrectDeletedEntries(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	_, _ = eng.Query(`rem @note doomed`)
	_, _ = eng.Query(`rem @note survivor`)

	matched, err := eng.Query(`forget contains["doomed"]`)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := eng.ProcessDeletion(nil, matched); err != nil {
		t.Fatal(err)
	}

	revised, err := eng.Query(`revise @important where all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(revised) != 1 {
		t.Fatalf("expected only the surviving entry to be revised, got %d", len(revised))
	}

	after, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(after) != 1 || after[0].Content != "survivor" {
		t.Fatalf("expected only the survivor after restart, got %v", after)
	}
}

func TestRecallByID(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @idea second`)
	_, _ = eng.Query(`rem @question third`)

	cases := []struct {
		query string
		ids   []int
	}{
		{`recall id[2]`, []int{2}},
		{`recall id[1,3]`, []int{1, 3}},
		{`recall id[3,1]`, []int{1, 3}},
		{`recall id[2,2]`, []int{2}},
		{`recall id[99]`, []int{}},
		{`recall id[1,2,3] AND type[idea]`, []int{2}},
		{`recall (id[1] OR id[3]) AND type[question]`, []int{3}},
		{`recall id[1] OR type[question]`, []int{1, 3}},
	}

	for _, c := range cases {
		results, err := eng.Query(c.query)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", c.query, err)
		}

		if len(results) != len(c.ids) {
			t.Fatalf("%s: expected ids %v, got %d results", c.query, c.ids, len(results))
		}

		for i, id := range c.ids {
			if results[i].ID != id {
				t.Fatalf("%s: expected ids %v, got a result with id %d", c.query, c.ids, results[i].ID)
			}
		}
	}
}

func TestForgetByID(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @note second`)
	_, _ = eng.Query(`rem @note third`)

	matched, err := eng.Query(`forget id[1,3]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(matched) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matched))
	}

	deleted, err := eng.ProcessDeletion(nil, matched)
	if err != nil {
		t.Fatal(err)
	}

	if deleted != 2 {
		t.Fatalf("expected 2 deletions, got %d", deleted)
	}

	remaining, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(remaining) != 1 || remaining[0].ID != 2 {
		t.Fatalf("expected only id 2 to survive, got %v", remaining)
	}
}

func TestReviseByID(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "chronicle.log")

	eng := openEngine(t, logPath)

	_, _ = eng.Query(`rem @note first`)
	_, _ = eng.Query(`rem @note second`)
	_, _ = eng.Query(`rem @note third`)

	results, err := eng.Query(`revise @important #flagged where id[1,3]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 entries revised, got %d", len(results))
	}

	after, err := openEngine(t, logPath).Query(`recall all`)
	if err != nil {
		t.Fatal(err)
	}

	if len(after) != 3 {
		t.Fatalf("expected 3 entries after restart, got %d", len(after))
	}

	for _, e := range after {
		if e.ID == 2 {
			if e.Type != entry.TypeNote || len(e.Tags) != 0 {
				t.Fatalf("id 2 should have been left alone, got %q %v", e.Type, e.Tags)
			}
			continue
		}

		if e.Type != entry.TypeImportant {
			t.Fatalf("id %d should be important, got %q", e.ID, e.Type)
		}

		if len(e.Tags) != 1 || e.Tags[0] != "flagged" {
			t.Fatalf("id %d should carry the new tag, got %v", e.ID, e.Tags)
		}
	}
}

func TestIDOfADeletedEntryMatchesNothing(t *testing.T) {
	eng := buildEngine(t)

	_, _ = eng.Query(`rem @note doomed`)
	_, _ = eng.Query(`rem @note survivor`)

	matched, _ := eng.Query(`forget id[1]`)
	if _, err := eng.ProcessDeletion(nil, matched); err != nil {
		t.Fatal(err)
	}

	results, err := eng.Query(`recall id[1]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatalf("expected a deleted id to match nothing, got %d", len(results))
	}
}
