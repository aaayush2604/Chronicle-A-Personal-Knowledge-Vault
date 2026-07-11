package terminal

import (
	"testing"
)

func completionsToStrings(r [][]rune) []string {
	out := make([]string, len(r))
	for i := range r {
		out[i] = string(r[i])
	}
	return out
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func assertContains(t *testing.T, got [][]rune, expected ...string) {
	t.Helper()

	completions := completionsToStrings(got)

	for _, e := range expected {
		if !contains(completions, e) {
			t.Fatalf("expected completion %q, got %v", e, completions)
		}
	}
}

func TestRecallSuggestions(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("recall "), len("recall "))

	if prefix != 0 {
		t.Fatalf("expected prefix length 0")
	}

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
		"time",
		"date",
		"len",
	)
}

func TestRecallPrefixCompletion(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("recall con"), len("recall con"))

	if prefix != 3 {
		t.Fatalf("expected prefix length 3 got %d", prefix)
	}

	assertContains(t, got, "tains[")
}

func TestRecallTypeSuggestions(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("recall type["), len("recall type["))

	if prefix != 0 {
		t.Fatalf("expected prefix length 0")
	}

	assertContains(
		t,
		got,
		"note",
		"learning",
		"question",
		"important",
		"idea",
	)
}

func TestRecallTypePrefixCompletion(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("recall type[n"), len("recall type[n"))

	if prefix != 1 {
		t.Fatalf("expected prefix length 1")
	}

	assertContains(t, got, "ote")
}

func TestRememberTypeSuggestions(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("rem @"), len("rem @"))

	if prefix != 0 {
		t.Fatalf("expected prefix length 0")
	}

	assertContains(
		t,
		got,
		"@note",
		"@learning",
		"@question",
		"@important",
		"@idea",
	)
}

func TestRememberTypePrefixCompletion(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("rem @q"), len("rem @q"))

	if prefix != 2 {
		t.Fatalf("expected prefix length 2")
	}

	assertContains(t, got, "uestion")
}

func TestAfterLogicalOperator(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall type[note] and "),
		len("recall type[note] and "),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestAfterOpeningParen(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall ("),
		len("recall ("),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestEmptyInputDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked on empty input: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune(""), 0)
}

func TestInvalidDoubleAtDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune("rem @@"), len("rem @@"))
}

func TestInvalidDoubleHashDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune("rem ##"), len("rem ##"))
}

func TestUnterminatedBracketDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune("recall type[["), len("recall type[["))
}

func TestUnterminatedStringDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune(`recall contains["hello`), len(`recall contains["hello`))
}

func TestRandomGarbageDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do([]rune("@#$%^&*"), len("@#$%^&*"))
}

func TestNestedExpressionDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall (type[note] and (tag[go"),
		len("recall (type[note] and (tag[go"),
	)
}

func TestRecallUnknownPrefix(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do([]rune("recall xyz"), len("recall xyz"))

	if len(got) != 0 {
		t.Fatalf("expected no completions")
	}
}

func TestRecallDoubleSpace(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do([]rune("recall  "), len("recall  "))

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
		"time",
		"date",
		"len",
	)
}

func TestRecallAfterOr(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall type[note] or "),
		len("recall type[note] or "),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestRememberFullCommand(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do([]rune("remember @"), len("remember @"))

	assertContains(
		t,
		got,
		"@note",
		"@learning",
		"@question",
		"@important",
		"@idea",
	)
}

func TestRememberLearningPrefix(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do([]rune("remember @l"), len("remember @l"))

	if prefix != 2 {
		t.Fatalf("expected prefix length 2 got %d", prefix)
	}

	assertContains(t, got, "earning")
}

func TestCursorInMiddleOfInput(t *testing.T) {
	c := NewCompleter()

	line := []rune("recall contains[")
	pos := len("recall ")

	got, _ := c.Do(line, pos)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestCursorInMiddleOfWord(t *testing.T) {
	c := NewCompleter()

	line := []rune("recall contains[")
	pos := len("recall con")

	got, prefix := c.Do(line, pos)

	if prefix != 3 {
		t.Fatalf("expected prefix length 3 got %d", prefix)
	}

	assertContains(t, got, "tains[")
}

func TestNestedParenthesesDoNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall (((type[n"),
		len("recall (((type[n"),
	)
}

func TestExtraClosingParenthesisDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall type[note]))"),
		len("recall type[note]))"),
	)
}

func TestRandomPunctuationDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("!@#$%^&*()"),
		len("!@#$%^&*()"),
	)
}

func TestRecallCaseInsensitive(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do([]rune("ReCaLl "), len("ReCaLl "))

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestNoDuplicateCompletions(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do([]rune("recall "), len("recall "))

	seen := make(map[string]struct{})

	for _, completion := range completionsToStrings(got) {
		if _, ok := seen[completion]; ok {
			t.Fatalf("duplicate completion %q", completion)
		}
		seen[completion] = struct{}{}
	}
}

func TestGroupingAfterOpenParen(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall ("),
		len("recall ("),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestGroupingAfterAnd(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (type[note] and "),
		len("recall (type[note] and "),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestGroupingAfterOr(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (type[note] or "),
		len("recall (type[note] or "),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestNestedGrouping(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (type[note] and ("),
		len("recall (type[note] and ("),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestGroupingTypeCompletion(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (type["),
		len("recall (type["),
	)

	assertContains(
		t,
		got,
		"note",
		"learning",
		"question",
		"important",
		"idea",
	)
}

func TestGroupingTagCompletionDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall (tags["),
		len("recall (tags["),
	)
}

func TestAfterClosedGroupDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall (type[note]) "),
		len("recall (type[note]) "),
	)
}

func TestGroupThenAnd(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (type[note]) and "),
		len("recall (type[note]) and "),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestDeepNestedGroupingDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall ((type[note] and (tags[go] or ("),
		len("recall ((type[note] and (tags[go] or ("),
	)
}

func TestUnclosedGroupingDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autocomplete panicked: %v", r)
		}
	}()

	c := NewCompleter()

	c.Do(
		[]rune("recall (type[note] and (tags[go]"),
		len("recall (type[note] and (tags[go]"),
	)
}

func TestMultipleNestedGroups(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall (((("),
		len("recall (((("),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}

func TestNestedGroupTypePrefixCompletion(t *testing.T) {
	c := NewCompleter()

	got, prefix := c.Do(
		[]rune("recall ((ty"),
		len("recall ((ty"),
	)

	if prefix != 2 {
		t.Fatalf("expected prefix length 2 got %d", prefix)
	}

	assertContains(t, got, "pe[")
}

func TestNestedGroupAfterLogicalOperator(t *testing.T) {
	c := NewCompleter()

	got, _ := c.Do(
		[]rune("recall ((type[note]) or ("),
		len("recall ((type[note]) or ("),
	)

	assertContains(
		t,
		got,
		"all",
		"contains[",
		"type[",
		"tags[",
	)
}
