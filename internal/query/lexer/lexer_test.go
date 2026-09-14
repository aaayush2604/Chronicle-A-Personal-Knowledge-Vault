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

func TestForgetKeyword(t *testing.T) {
	scanner := NewScanner("forget all")

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}

	if tokens[0].TokenType != COMMAND {
		t.Fatalf("expected COMMAND, got %v", tokens[0].TokenType)
	}

	if tokens[0].Lexeme != "forget" {
		t.Fatalf("expected forget, got %q", tokens[0].Lexeme)
	}

	if tokens[1].TokenType != ALL {
		t.Fatalf("expected ALL token, got %v", tokens[1].TokenType)
	}
}

func TestPunctuationTokens(t *testing.T) {
	tests := []string{
		",",
		".",
		"?",
		"'",
	}

	for _, input := range tests {
		scanner := NewScanner(input)

		tokens, err := scanner.ScanTokens()
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", input, err)
		}

		if tokens[0].TokenType != PUNCTUATION {
			t.Fatalf("%q should be PUNCTUATION", input)
		}
	}
}

func TestNotEqualOperator(t *testing.T) {
	scanner := NewScanner("len != 10")

	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tokens[1].TokenType != OPERATOR {
		t.Fatalf("expected OPERATOR, got %v", tokens[1].TokenType)
	}

	if tokens[1].Lexeme != "!=" {
		t.Fatalf("expected !=, got %q", tokens[1].Lexeme)
	}
}

func TestRemContentIsVerbatim(t *testing.T) {
	cases := []struct {
		input   string
		content string
	}{
		{`rem I will revise this tomorrow`, "I will revise this tomorrow"},
		{`rem where did I put my keys`, "where did I put my keys"},
		{`rem recall the meeting notes`, "recall the meeting notes"},
		{`rem all hands meeting`, "all hands meeting"},
		{`rem Meeting with John, at 5 PM.`, "Meeting with John, at 5 PM."},
		{`rem cost is 50% of budget`, "cost is 50% of budget"},
		{`rem see http://example.com and/or ask`, "see http://example.com and/or ask"},
		{`rem my e-mail id`, "my e-mail id"},
		{`rem café naïve 🎉`, "café naïve 🎉"},
		{`rem pipe | inside content`, "pipe | inside content"},
		{`remember note: this is a test`, "note: this is a test"},
		{`rem @idea #go #db all hands sync at 4:30`, "all hands sync at 4:30"},
		{`rem #go @idea the text`, "the text"},
	}

	for _, c := range cases {
		tokens, err := NewScanner(c.input).ScanTokens()
		if err != nil {
			t.Fatalf("%q: unexpected error %v", c.input, err)
		}

		var got string
		for _, tok := range tokens {
			if tok.TokenType == TEXT {
				got = tok.Lexeme
			}
		}

		if got != c.content {
			t.Fatalf("%q: expected content %q, got %q", c.input, c.content, got)
		}
	}
}

func TestRemPrefixIsStillTokenized(t *testing.T) {
	tokens, err := NewScanner(`rem @idea #go #db the content`).ScanTokens()
	if err != nil {
		t.Fatal(err)
	}

	expected := []TokenType{COMMAND, ETYPE, TAG, TAG, TEXT, EOF}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, tType := range expected {
		if tokens[i].TokenType != tType {
			t.Fatalf("token %d: expected %v, got %v", i, tType, tokens[i].TokenType)
		}
	}
}

func TestRemWithoutContentHasNoTextToken(t *testing.T) {
	for _, input := range []string{"rem", "rem ", "rem @idea", "rem @idea #go "} {
		tokens, err := NewScanner(input).ScanTokens()
		if err != nil {
			t.Fatalf("%q: unexpected error %v", input, err)
		}

		for _, tok := range tokens {
			if tok.TokenType == TEXT {
				t.Fatalf("%q: expected no content token, got %q", input, tok.Lexeme)
			}
		}
	}
}

func TestQueryCommandsStillValidateCharacters(t *testing.T) {
	for _, input := range []string{`recall my e-mail`, `forget 50%`, `revise @note where x:y`} {
		if _, err := NewScanner(input).ScanTokens(); err == nil {
			t.Fatalf("%q: expected a scan error, query syntax should stay strict", input)
		}
	}
}
