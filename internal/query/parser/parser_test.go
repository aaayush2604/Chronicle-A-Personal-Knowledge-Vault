package parser

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
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
