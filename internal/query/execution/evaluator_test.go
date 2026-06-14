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
	tokens, err := scanner.ScanTokens()

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

	if !EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
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

	if EvaluateRecall(expr, record) {
		t.Fatalf("expected false")
	}
}

func TestEvaluateTypeFilterSuccess(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[note]`,
	)

	record := Record{
		"type": "note",
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateTypeFilterFailure(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[note]`,
	)

	record := Record{
		"type": "idea",
	}

	if EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
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

	if EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
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

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestTypeAliasNote(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[n]`,
	)

	record := Record{
		"type": "note",
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected alias n to match note")
	}
}

func TestTypeAliasIdea(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[i]`,
	)

	record := Record{
		"type": "idea",
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected alias i to match idea")
	}
}

func TestTypeAliasQuestion(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[q]`,
	)

	record := Record{
		"type": "question",
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected alias q to match question")
	}
}

func TestTypeAliasLearning(t *testing.T) {
	expr := parseExpr(
		t,
		`recall type[l]`,
	)

	record := Record{
		"type": "learning",
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected alias l to match learning")
	}
}

func TestEvaluateTagsSuccess(t *testing.T) {
	expr := parseExpr(
		t,
		`recall tags[go,database]`,
	)

	record := Record{
		"tags": []string{
			"database",
			"project",
		},
	}

	if !EvaluateRecall(expr, record) {
		t.Fatalf("expected true")
	}
}

func TestEvaluateTagsFailure(t *testing.T) {
	expr := parseExpr(
		t,
		`recall tags[go,database]`,
	)

	record := Record{
		"tags": []string{
			"python",
			"backend",
		},
	}

	if EvaluateRecall(expr, record) {
		t.Fatalf("expected false")
	}
}
