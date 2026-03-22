package parser

import (
	"chronicle/internal/errorC"
	"fmt"
)

type Parser struct {
	Tokens  []*Token
	Current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) matchLexeme(lexemes ...string) bool {
	for _, lexeme := range lexemes {
		if p.check(lexeme) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) checkTokenType(tTypes ...TokenType) bool {
	for _, tType := range tTypes {
		if p.peek().TokenType == tType {
			return true
		}
	}
	return false
}

func (p *Parser) check(lexeme string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Lexeme == lexeme
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() *Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() *Token {
	return p.Tokens[p.Current-1]
}

// of no use yet, since only one expression tree is allowed
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.check("recall") {
			return
		}

		if p.check("AND") || p.check("OR") {
			p.advance()
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() (*Query, error) {
	q, err := p.parseQuery()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Query")
	}
	return q, nil
}

// Now we start adding a function to parse each of the grammar rules
func (p *Parser) parseQuery() (*Query, error) {

	if !p.matchLexeme("recall") {
		err := errorC.New(errorC.Syntax, "Syntax Error: Query must begin with a command")
		return nil, err
	}

	expr, err := p.parseExpression()

	if err != nil {
		err := errorC.Wrap(err, errorC.Syntax, "Error parsing Expression:")
		return nil, err
	}

	return &Query{
		Command: RecallCommand,
		Expr:    expr,
	}, nil
}

func (p *Parser) parseExpression() (Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Term:")
	}

	for p.matchLexeme("or") {
		OrToken := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Term:")
		}
		expr = NewLogical(expr, OrToken, right)
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Factor:")
	}

	for p.matchLexeme("and") {
		AndToken := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Factor:")
		}
		expr = NewLogical(expr, AndToken, right)
	}

	return expr, nil
}

func (p *Parser) parseFactor() (Expr, error) {
	if p.checkTokenType(LPAREN) {
		expr, err := p.parseGrouping()
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Grouping:")
		}
		return expr, nil
	} else if p.checkTokenType(IDENTIFIER) {
		if p.matchLexeme("contains") {
			expr, err := p.parseContains()
			if err != nil {
				return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Contains:")
			}
			return expr, nil
		} else if p.matchLexeme("type") {
			expr, err := p.parseTypeFilter()
			if err != nil {
				return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Type Filter:")
			}
			return expr, nil
		} else {
			expr, err := p.parseComparison()
			if err != nil {
				return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Comparison:")
			}
			return expr, nil
		}
	}
	err := errorC.New(errorC.Syntax, fmt.Sprintf("Invalid Token at col %d", p.peek().Position))
	return nil, err
}

func (p *Parser) parseGrouping() (Expr, error) {
	if !p.matchLexeme("(") {
		return nil, errorC.New(errorC.Internal, fmt.Sprintf("Issue in Lexer. LPAREN Token with wrong Lexeme Found at col %d", p.peek().Position))
	}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing Expression under Grouping")
	}
	if !p.matchLexeme(")") {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Grouping at col %d", p.peek().Position))
	}
	expr = NewGrouping(expr)
	return expr, nil
}

func (p *Parser) parseContains() (Expr, error) {

	if !(p.checkTokenType(LBRACKET) && p.matchLexeme("[")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Contains StringList [ at col %d", p.peek().Position))
	}

	strings, err := p.parseStringList()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing StringList:")
	}
	expr := NewContains(strings)

	if !(p.checkTokenType(RBRACKET) && p.matchLexeme("]")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Contains StringList ] at col %d", p.peek().Position))
	}

	return expr, nil
}

func (p *Parser) parseTypeFilter() (Expr, error) {

	if !(p.checkTokenType(LBRACKET) && p.matchLexeme("[")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing TypeFilter StringList [ at col %d", p.peek().Position))
	}

	strings, err := p.parseStringList()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing StringList:")
	}
	expr := NewTypeFilter(strings)

	if !(p.checkTokenType(RBRACKET) && p.matchLexeme("]")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing TypeFilter StringList ] at col %d", p.peek().Position))
	}

	return expr, nil
}

func (p *Parser) parseStringList() ([]string, error) {
	var words []string
	if !p.checkTokenType(STRING) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a String at col %d", p.peek().Position))
	}

	words = append(words, p.peek().Literal.(string))
	p.advance()

	for p.checkTokenType(COMMA) {
		p.advance()
		if !p.checkTokenType(STRING) {
			return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a String or ] at col %d", p.peek().Position))
		}
		words = append(words, p.peek().Literal.(string))
		p.advance()
	}
	return words, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	if !p.checkTokenType(IDENTIFIER) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Identifier at col: %d", p.peek().Position))
	}
	field := p.peek()
	p.advance()

	if !p.checkTokenType(OPERATOR) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Operator at col: %d", p.peek().Position))
	}
	opToken := p.peek()
	p.advance()

	if !p.checkTokenType(STRING, NUMBER) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a Literal(Number or String) at col: %d", p.peek().Position))
	}
	literal := *NewLiteral(p.peek().Literal)
	p.advance()

	return NewComparison(field, opToken, literal), nil
}
