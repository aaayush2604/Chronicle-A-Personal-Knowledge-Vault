package parser

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
	"strings"
	"testing"
)

func parseQuery(t *testing.T, input string) *Query {
	t.Helper()

	scanner := lexer.NewScanner(input)
	tokens, err := scanner.ScanTokens()

	p := NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	return q
}

func TestParseAllQuery(t *testing.T) {
	q := parseQuery(t, `recall all`)

	if q.Command != RecallCommand {
		t.Fatalf("expected recall command")
	}

	if _, ok := q.Expr.(*All); !ok {
		t.Fatalf("expected All expression")
	}

	q = parseQuery(t, `forget all`)

	if q.Command != ForgetCommand {
		t.Fatalf("expected forget command")
	}

	if _, ok := q.Expr.(*All); !ok {
		t.Fatalf("expected All expression")
	}
}

func TestParseContains(t *testing.T) {
	q := parseQuery(
		t,
		`recall contains["go","database"]`,
	)

	expr, ok := q.Expr.(*Contains)
	if !ok {
		t.Fatalf("expected Contains expression")
	}

	if len(expr.Words) != 2 {
		t.Fatalf("expected 2 words")
	}
}

func TestParseTypeFilter(t *testing.T) {
	q := parseQuery(
		t,
		`recall type[note,idea]`,
	)

	expr, ok := q.Expr.(*TypeFilter)
	if !ok {
		t.Fatalf("expected TypeFilter expression")
	}

	if len(expr.Words) != 2 {
		t.Fatalf("expected 2 words")
	}
}

func TestParseComparison(t *testing.T) {
	q := parseQuery(
		t,
		`recall len > 10`,
	)

	expr, ok := q.Expr.(*Comparison)
	if !ok {
		t.Fatalf("expected Comparison")
	}

	if expr.Field.Lexeme != "len" {
		t.Fatalf("unexpected field")
	}

	if expr.Operator.Lexeme != ">" {
		t.Fatalf("unexpected operator")
	}
}

func TestParseLogicalAnd(t *testing.T) {
	q := parseQuery(
		t,
		`recall len > 10 and len < 100`,
	)

	_, ok := q.Expr.(*Logical)

	if !ok {
		t.Fatalf("expected Logical expression")
	}
}

func TestParseLogicalOr(t *testing.T) {
	q := parseQuery(
		t,
		`recall len > 10 or len < 100`,
	)

	_, ok := q.Expr.(*Logical)

	if !ok {
		t.Fatalf("expected Logical expression")
	}
}

func TestParseGrouping(t *testing.T) {
	q := parseQuery(
		t,
		`recall (len > 10)`,
	)

	_, ok := q.Expr.(*Grouping)

	if !ok {
		t.Fatalf("expected Grouping")
	}
}

func TestParseMissingRecall(t *testing.T) {
	scanner := lexer.NewScanner(`len > 10`)
	tokens, err := scanner.ScanTokens()

	p := NewParser(tokens)

	_, err = p.Parse()

	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestParseInvalidGrouping(t *testing.T) {
	scanner := lexer.NewScanner(
		`recall (len > 10`,
	)

	tokens, err := scanner.ScanTokens()

	p := NewParser(tokens)

	_, err = p.Parse()

	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestParseContainsWithoutBracket(t *testing.T) {
	scanner := lexer.NewScanner(
		`recall contains "go"`,
	)

	tokens, err := scanner.ScanTokens()

	p := NewParser(tokens)

	_, err = p.Parse()

	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestParseComparisonWithoutValue(t *testing.T) {
	scanner := lexer.NewScanner(
		`recall len >`,
	)

	tokens, err := scanner.ScanTokens()

	p := NewParser(tokens)

	_, err = p.Parse()

	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestAndHasHigherPrecedenceThanOr(t *testing.T) {
	q := parseQuery(
		t,
		`recall len > 100 or len > 5 and len < 50`,
	)

	root, ok := q.Expr.(*Logical)
	if !ok {
		t.Fatalf("expected root logical node")
	}

	if root.Operator.Lexeme != "or" {
		t.Fatalf(
			"expected OR at root, got %s",
			root.Operator.Lexeme,
		)
	}
}

func TestParseRemCommand(t *testing.T) {
	q := parseQuery(
		t,
		`rem #go #database "hello world"`,
	)

	if q.Command != RemCommand {
		t.Fatalf("expected note command")
	}
}

func TestParseTags(t *testing.T) {
	q := parseQuery(
		t,
		`recall tags[go,database]`,
	)

	expr, ok := q.Expr.(*Tags)

	if !ok {
		t.Fatalf("expected Tags expression")
	}

	if len(expr.TagList) != 2 {
		t.Fatalf("expected 2 tags")
	}
}

func TestParseEmptyTags(t *testing.T) {
	q := parseQuery(
		t,
		`recall tags[]`,
	)

	_, ok := q.Expr.(*Tags)

	if !ok {
		t.Fatalf("expected Tags expression")
	}
}

func TestParseRememberCommand(t *testing.T) {
	q := parseQuery(
		t,
		`remember #go hello`,
	)

	if q.Command != RemCommand {
		t.Fatalf("expected remember command")
	}
}

func TestParseRemDefaultType(t *testing.T) {
	q := parseQuery(
		t,
		`rem hello world`,
	)

	payload := q.Payload.(*RemPayload)

	if payload.Type != entry.TypeNote {
		t.Fatalf("expected default note type")
	}
}

func TestParseForgetAll(t *testing.T) {
	scanner := lexer.NewScanner("forget all")

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if q.Command != ForgetCommand {
		t.Fatalf("expected ForgetCommand, got %v", q.Command)
	}

	if _, ok := q.Expr.(*All); !ok {
		t.Fatalf("expected *All expression")
	}
}

func TestParseForgetContains(t *testing.T) {
	scanner := lexer.NewScanner(`forget contains["database"]`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if q.Command != ForgetCommand {
		t.Fatalf("expected ForgetCommand")
	}

	if _, ok := q.Expr.(*Contains); !ok {
		t.Fatalf("expected Contains expression")
	}
}

func TestParseForgetLogicalExpression(t *testing.T) {
	scanner := lexer.NewScanner(
		`forget type[note] AND contains["go"]`,
	)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := q.Expr.(*Logical); !ok {
		t.Fatalf("expected Logical expression")
	}
}

func TestParseForgetComparison(t *testing.T) {
	scanner := lexer.NewScanner(
		`forget len != 100`,
	)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	cmp, ok := q.Expr.(*Comparison)
	if !ok {
		t.Fatalf("expected Comparison")
	}

	if cmp.Operator.Lexeme != "!=" {
		t.Fatalf("expected != operator")
	}
}

func TestRemPayloadAcceptsEitherOrder(t *testing.T) {
	cases := []struct {
		input   string
		eType   entry.EntryType
		tags    []string
		content string
	}{
		{`rem the text`, entry.TypeNote, nil, "the text"},
		{`rem @idea the text`, entry.TypeIdea, nil, "the text"},
		{`rem @idea #go the text`, entry.TypeIdea, []string{"go"}, "the text"},
		{`rem #go @idea the text`, entry.TypeIdea, []string{"go"}, "the text"},
		{`rem #a @idea #b the text`, entry.TypeIdea, []string{"a", "b"}, "the text"},
		{`rem #a #b @idea the text`, entry.TypeIdea, []string{"a", "b"}, "the text"},
		{`rem @q shorthand alias`, entry.TypeQuestion, nil, "shorthand alias"},
		{`rem #go #go deduped`, entry.TypeNote, []string{"go"}, "deduped"},
		{`rem #my-tag hyphenated`, entry.TypeNote, []string{"my-tag"}, "hyphenated"},
		{`rem # bare hash`, entry.TypeNote, nil, "bare hash"},
		{`rem @idea "Read chapter 3"`, entry.TypeIdea, nil, "Read chapter 3"},
		{`rem "@todo buy milk"`, entry.TypeNote, nil, "@todo buy milk"},
		{`rem say "hello" to him`, entry.TypeNote, nil, `say "hello" to him`},
		{`rem email john@example.com`, entry.TypeNote, nil, "email john@example.com"},
		{`rem issue #42 stays inline`, entry.TypeNote, nil, "issue #42 stays inline"},
	}

	for _, c := range cases {
		q := parseQuery(t, c.input)

		payload, ok := q.Payload.(*RemPayload)
		if !ok {
			t.Fatalf("%q: expected a rem payload", c.input)
		}

		if payload.Type != c.eType {
			t.Fatalf("%q: expected type %q, got %q", c.input, c.eType, payload.Type)
		}

		if payload.Content != c.content {
			t.Fatalf("%q: expected content %q, got %q", c.input, c.content, payload.Content)
		}

		if len(payload.Tags) != len(c.tags) {
			t.Fatalf("%q: expected tags %v, got %v", c.input, c.tags, payload.Tags)
		}

		for i := range c.tags {
			if payload.Tags[i] != c.tags[i] {
				t.Fatalf("%q: expected tags %v, got %v", c.input, c.tags, payload.Tags)
			}
		}
	}
}

func TestRemPayloadRejectsUnknownType(t *testing.T) {
	tokens, err := lexer.NewScanner(`rem @todo buy milk`).ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatalf("expected an error for an unknown type")
	}

	for _, want := range []string{"@todo", "@note", "@idea", "@question", "@learning", "@important", "quote"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected the error to mention %q, got %q", want, err.Error())
		}
	}
}

func TestRemPayloadRejectsTwoTypes(t *testing.T) {
	tokens, err := lexer.NewScanner(`rem @idea @question the text`).ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if _, err = NewParser(tokens).Parse(); err == nil {
		t.Fatalf("expected an error when a type is given twice")
	}
}

func TestTrailingInputIsRejected(t *testing.T) {
	cases := []string{
		`recall type[note] this is garbage`,
		`recall tags[go] extra`,
		`forget type[note] junk`,
		`revise @idea where type[note] junk`,
		`recall len > 10 leftover`,
	}

	for _, input := range cases {
		tokens, err := lexer.NewScanner(input).ScanTokens()
		if err != nil {
			t.Fatalf("%q: unexpected scan error %v", input, err)
		}

		if _, err := NewParser(tokens).Parse(); err == nil {
			t.Fatalf("%q: expected trailing input to be rejected", input)
		}
	}
}

func TestCompleteQueriesAreStillAccepted(t *testing.T) {
	cases := []string{
		`recall all`,
		`recall type[note]`,
		`recall type[note] AND tags[go]`,
		`recall (type[note] OR type[idea]) AND len > 10`,
		`forget contains["go"]`,
		`revise @idea where type[note]`,
		`revise @idea where all`,
	}

	for _, input := range cases {
		parseQuery(t, input)
	}
}

func TestReservedWordsInListsExplainThemselves(t *testing.T) {
	cases := []struct {
		input string
		word  string
		list  string
	}{
		{`recall tags[all]`, "all", "tags"},
		{`recall tags[where]`, "where", "tags"},
		{`recall tags[go,or]`, "or", "tags"},
		{`recall type[and]`, "and", "type"},
	}

	for _, c := range cases {
		tokens, err := lexer.NewScanner(c.input).ScanTokens()
		if err != nil {
			t.Fatalf("%q: unexpected scan error %v", c.input, err)
		}

		_, err = NewParser(tokens).Parse()
		if err == nil {
			t.Fatalf("%q: expected an error", c.input)
		}

		message := err.Error()
		for _, want := range []string{c.word, "reserved word", c.list + "[...]", "help"} {
			if !strings.Contains(message, want) {
				t.Fatalf("%q: expected the error to mention %q, got %q", c.input, want, message)
			}
		}

		for _, reserved := range lexer.ReservedWords() {
			if !strings.Contains(message, reserved) {
				t.Fatalf("%q: expected the error to list %q, got %q", c.input, reserved, message)
			}
		}
	}
}

func TestUnknownWordInListKeepsTheSimpleError(t *testing.T) {
	tokens, err := lexer.NewScanner(`recall tags[123]`).ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatalf("expected an error")
	}

	if strings.Contains(err.Error(), "reserved word") {
		t.Fatalf("a non keyword should not be reported as reserved, got %q", err.Error())
	}
}

func TestParseIDList(t *testing.T) {
	cases := []struct {
		input string
		ids   []int
	}{
		{`recall id[3]`, []int{3}},
		{`recall id[3,7,12]`, []int{3, 7, 12}},
		{`recall id[3, 7]`, []int{3, 7}},
		{`recall id[3,3,7]`, []int{3, 7}},
		{`forget id[1]`, []int{1}},
	}

	for _, c := range cases {
		q := parseQuery(t, c.input)

		ids, ok := q.Expr.(*IDs)
		if !ok {
			t.Fatalf("%q: expected an id predicate", c.input)
		}

		if len(ids.List) != len(c.ids) {
			t.Fatalf("%q: expected %v, got %v", c.input, c.ids, ids.List)
		}

		for i := range c.ids {
			if ids.List[i] != c.ids[i] {
				t.Fatalf("%q: expected %v, got %v", c.input, c.ids, ids.List)
			}
		}
	}
}

func TestParseIDListRejectsNonIDs(t *testing.T) {
	for _, input := range []string{
		`recall id[abc]`,
		`recall id[1.5]`,
		`recall id[0]`,
		`recall id[1,]`,
		`recall id[1 2]`,
		`recall id[1`,
		`recall id`,
		`recall id["3"]`,
	} {
		tokens, err := lexer.NewScanner(input).ScanTokens()
		if err != nil {
			continue
		}

		if _, err := NewParser(tokens).Parse(); err == nil {
			t.Fatalf("%q: expected an error", input)
		}
	}
}

func TestParseIDListInsideExpressions(t *testing.T) {
	for _, input := range []string{
		`recall id[3] AND type[note]`,
		`recall (id[1] OR id[2]) AND type[idea]`,
		`forget id[3,7] OR tags[draft]`,
		`revise @important where id[3,7]`,
	} {
		parseQuery(t, input)
	}
}
