package execution

import (
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"testing"
	"time"
)

func parseExpr(t *testing.T, input string) parser.Expr {
	t.Helper()

	scanner := lexer.NewScanner(input)
	tokens := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	return q.Expr
}

func TestEvaluateAll(t *testing.T) {
	expr := parseExpr(t, `recall all`)

	record := Record{
		"content": "hello",
		"len":     5,
		"type":    "note",
		"date":    time.Now(),
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateContainsSuccess(t *testing.T) {
	expr := parseExpr(
		t,
		`recall contains["database"]`,
	)

	record := Record{
		"content": "I love database systems",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateContainsFailure(t *testing.T) {
	expr := parseExpr(
		t,
		`recall contains["database"]`,
	)

	record := Record{
		"content": "golang programming",
	}

	if Evaluate(expr, record) {
		t.Fatalf("expected false")
	}
}

func TestEvaluateTypeFilterSuccess(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["note"]`,
	)

	record := Record{
		"type": "note",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateTypeFilterFailure(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["note"]`,
	)

	record := Record{
		"type": "idea",
	}

	if Evaluate(expr, record) {
		t.Fatalf("expected false")
	}
}

func TestEvaluateLenGreaterThan(t *testing.T) {
	expr := parseExpr(
		t,
		`recall len > 10`,
	)

	record := Record{
		"len": 20,
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateLenLessThan(t *testing.T) {
	expr := parseExpr(
		t,
		`recall len < 10`,
	)

	record := Record{
		"len": 20,
	}

	if Evaluate(expr, record) {
		t.Fatalf("expected false")
	}
}

func TestEvaluateAndExpression(t *testing.T) {
	expr := parseExpr(
		t,
		`recall len > 10 and len < 30`,
	)

	record := Record{
		"len": 20,
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateOrExpression(t *testing.T) {
	expr := parseExpr(
		t,
		`recall len > 100 or len < 30`,
	)

	record := Record{
		"len": 20,
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateDateComparison(t *testing.T) {
	expr := parseExpr(
		t,
		`recall date >= "2025"`,
	)

	record := Record{
		"date": time.Date(
			2026,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateTimeComparison(t *testing.T) {
	expr := parseExpr(
		t,
		`recall time > "7 PM"`,
	)

	record := Record{
		"date": time.Date(
			2026,
			1,
			1,
			21,
			0,
			0,
			0,
			time.UTC,
		),
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestTypeAliasNote(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["n"]`,
	)

	record := Record{
		"type": "note",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected alias n to match note")
	}
}

func TestTypeAliasIdea(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["i"]`,
	)

	record := Record{
		"type": "idea",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected alias i to match idea")
	}
}

func TestTypeAliasQuestion(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["q"]`,
	)

	record := Record{
		"type": "question",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected alias q to match question")
	}
}

func TestTypeAliasLearning(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type["l"]`,
	)

	record := Record{
		"type": "learning",
	}

	if !Evaluate(expr, record) {
		t.Fatalf("expected alias l to match learning")
	}
}
