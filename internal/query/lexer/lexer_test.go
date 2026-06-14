package lexer

import "testing"

func TestScanRecallCommand(t *testing.T) {
	scanner := NewScanner(`recall all`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}

	if tokens[0].TokenType != COMMAND {
		t.Fatalf("expected COMMAND token")
	}

	if tokens[1].TokenType != ALL {
		t.Fatalf("expected ALL token")
	}

	if tokens[2].TokenType != EOF {
		t.Fatalf("expected EOF token")
	}
}

func TestScanStringLiteral(t *testing.T) {
	scanner := NewScanner(`"hello world"`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if tokens[0].TokenType != STRING {
		t.Fatalf("expected STRING token")
	}

	if tokens[0].Literal != "hello world" {
		t.Fatalf("expected literal hello world")
	}
}

func TestScanNumberLiteral(t *testing.T) {
	scanner := NewScanner(`123`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if tokens[0].TokenType != NUMBER {
		t.Fatalf("expected NUMBER token")
	}

	if tokens[0].Literal.(float64) != 123 {
		t.Fatalf("expected 123")
	}
}

func TestScanDecimalNumber(t *testing.T) {
	scanner := NewScanner(`12.5`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if tokens[0].TokenType != NUMBER {
		t.Fatalf("expected NUMBER token")
	}

	if tokens[0].Literal.(float64) != 12.5 {
		t.Fatalf("expected 12.5")
	}
}

func TestScanLogicalOperators(t *testing.T) {
	scanner := NewScanner(`AND OR`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if tokens[0].TokenType != LOGICAL {
		t.Fatalf("expected LOGICAL")
	}

	if tokens[1].TokenType != LOGICAL {
		t.Fatalf("expected LOGICAL")
	}
}

func TestScanComparisonOperators(t *testing.T) {
	tests := []string{
		"<",
		">",
		"<=",
		">=",
		"=",
		"!=",
	}

	for _, input := range tests {
		scanner := NewScanner(input)

		tokens, err := scanner.ScanTokens()
		if err != nil {
			t.Fatalf("%s", err.Error())
		}

		if tokens[0].TokenType != OPERATOR {
			t.Fatalf("expected OPERATOR for %s", input)
		}
	}
}

func TestScanBrackets(t *testing.T) {
	scanner := NewScanner(`[]()`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	expected := []TokenType{
		LBRACKET,
		RBRACKET,
		LPAREN,
		RPAREN,
		EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("unexpected token count")
	}

	for i := range expected {
		if tokens[i].TokenType != expected[i] {
			t.Fatalf(
				"expected %v at position %d got %v",
				expected[i],
				i,
				tokens[i].TokenType,
			)
		}
	}
}

func TestScanIdentifier(t *testing.T) {
	scanner := NewScanner(`contains`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if tokens[0].TokenType != IDENTIFIER {
		t.Fatalf("expected IDENTIFIER")
	}
}

func TestScanCaseInsensitiveKeywords(t *testing.T) {
	scanner := NewScanner(`ReCaLl AnD Or AlL`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	expected := []TokenType{
		COMMAND,
		LOGICAL,
		LOGICAL,
		ALL,
		EOF,
	}

	for i := range expected {
		if tokens[i].TokenType != expected[i] {
			t.Fatalf(
				"expected %v got %v",
				expected[i],
				tokens[i].TokenType,
			)
		}
	}
}

func TestScanComplexQuery(t *testing.T) {
	scanner := NewScanner(
		`recall date >= "2025" AND type["note"]`,
	)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if len(tokens) == 0 {
		t.Fatalf("expected tokens")
	}

	if tokens[len(tokens)-1].TokenType != EOF {
		t.Fatalf("expected EOF")
	}
}

func TestScanRemCommand(t *testing.T) {
	scanner := NewScanner(`rem hello`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].TokenType != COMMAND {
		t.Fatalf("expected COMMAND")
	}
}

func TestScanRememberCommand(t *testing.T) {
	scanner := NewScanner(`remember hello`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].TokenType != COMMAND {
		t.Fatalf("expected COMMAND")
	}
}

func TestScanNoteEntryType(t *testing.T) {
	scanner := NewScanner(`note`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].TokenType != ETYPE {
		t.Fatalf("expected ETYPE")
	}
}

func TestScanNoteAlias(t *testing.T) {
	scanner := NewScanner(`n`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].TokenType != ETYPE {
		t.Fatalf("expected ETYPE")
	}
}

func TestScanTag(t *testing.T) {
	scanner := NewScanner(`#golang`)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].TokenType != TAG {
		t.Fatalf("expected TAG")
	}

	if tokens[0].Literal.(string) != "golang" {
		t.Fatalf("expected golang")
	}
}

func TestInvalidDoubleHashTag(t *testing.T) {
	scanner := NewScanner(`##go`)

	_, err := scanner.ScanTokens()

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestUnterminatedString(t *testing.T) {
	scanner := NewScanner(`"hello`)

	_, err := scanner.ScanTokens()

	if err == nil {
		t.Fatalf("expected error")
	}
}
