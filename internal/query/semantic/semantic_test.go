package semantic

import (
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"testing"
)

func parseForSemantic(t *testing.T, input string) *parser.Query {
	t.Helper()

	scanner := lexer.NewScanner(input)
	tokens, err := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	return q
}

func TestValidLenComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall len > 10`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidDateComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall date >= "2025"`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidTimeComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall time > "7 PM"`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidContains(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall contains["go"]`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidTypeFilter(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall type[note]`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidAll(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall all`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvalidField(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall foo > 10`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestInvalidDateComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall date > 10`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestInvalidTimeComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall time > 10`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestInvalidLenComparison(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall len > "abc"`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestEmptyContains(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall contains[]`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestEmptyTypeFilter(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall type[]`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestValidLogicalExpression(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall len > 10 and date > "2025"`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidGroupedExpression(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall (len > 10 and date > "2025")`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmptyTags(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall tags[]`,
	)

	err := AnalyzeSemantics(q)

	if err == nil {
		t.Fatalf("expected semantic error")
	}
}

func TestValidTags(t *testing.T) {
	q := parseForSemantic(
		t,
		`recall tags[go]`,
	)

	if err := AnalyzeSemantics(q); err != nil {
		t.Fatalf("unexpected error")
	}
}
