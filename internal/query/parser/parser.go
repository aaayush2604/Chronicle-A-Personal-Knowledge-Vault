package parser

import (
	"chronicle/internal/errorC"
	"fmt"
)

type Parser struct {
	Tokens  []Token
	Current int
}

func NewParser(tokens []Token) *Parser {
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

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

//Now we start adding a function to parse each of the grammar rules

func (p *Parser) parseQuery() (*Query, error) {
	var expr Expr
	var err error = nil
	if p.matchLexeme("recall") {
		expr, err = p.parseExpression()
		return &Query{
			Command: RecallCommand,
			Expr:    expr,
		}, nil
	}
	err = errorC.New(errorC.NotFound, "Query Command Not Found, Usage: recall <clauses>")
	return nil, err
}

func (p *Parser) parseExpression() (Expr, error) {
	var expr Expr
	var err error = nil
	var right Expr

	expr, err = p.parseTerm()

	for p.matchLexeme("OR") {
		OrToken := p.previous()
		right, err = p.parseTerm()
		expr = NewLogical(expr, OrToken, right)
	}

	return expr, err
}

func (p *Parser) parseTerm() (Expr, error) {
	var expr Expr
	var err error = nil
	var right Expr

	expr, err = p.parseFactor()

	for p.matchLexeme("AND") {
		AndToken := p.previous()
		right, err = p.parseFactor()
		expr = NewLogical(expr, AndToken, right)
	}

	return expr, err
}

func (p *Parser) parseFactor() (Expr, error) {
	var expr Expr
	var err error = nil
	if p.checkTokenType(LPAREN) {
		expr, err = p.parseGrouping()
	} else if p.checkTokenType(IDENTIFIER) {
		if p.check("contains") {
			expr, err = p.parseContains()
		} else if p.check("type") {
			expr, err = p.parseTypeFilter()
		} else {
			expr, err = p.parseComparison()
		}
	} else {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Invalid Token at col %d", p.peek().Position))
		return nil, err
	}
	return expr, err
}

func (p *Parser) parseGrouping() (Expr, error) {
	var err error = nil
	if !p.matchLexeme("(") {
		err = errorC.New(errorC.Internal, fmt.Sprintf("Issue in Lexer. LPAREN Token with wrong Lexeme Found at col %d", p.peek().Position))
	}
	expr, err := p.parseExpression()
	if !p.matchLexeme(")") {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse Grouping at col %d", p.peek().Position))
	}
	expr = NewGrouping(expr)
	return expr, err
}

func (p *Parser) parseContains() (Expr, error) {
	var err error = nil
	var expr Expr
	//consume the contains token
	p.advance()

	if !(p.checkTokenType(LBRACKET) && p.matchLexeme("[")) {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse StringList at col %d", p.peek().Position))
	}

	strings, err := p.parseStringList()
	expr = NewContains(strings)

	if !(p.checkTokenType(RBRACKET) && p.matchLexeme("]")) {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse StringList at col %d", p.peek().Position))
	}

	return expr, err
}

func (p *Parser) parseTypeFilter() (Expr, error) {
	var err error = nil
	var expr Expr
	//consume the contains token
	p.advance()

	if !(p.checkTokenType(LBRACKET) && p.matchLexeme("[")) {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse StringList at col %d", p.peek().Position))
	}

	strings, err := p.parseStringList()
	expr = NewTypeFilter(strings)

	if !(p.checkTokenType(RBRACKET) && p.matchLexeme("]")) {
		err = errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse StringList at col %d", p.peek().Position))
	}

	return expr, err
}

func (p *Parser) parseStringList() ([]string, error) {
	var err error = nil
	var words []string
	if p.checkTokenType(STRING) {
		words = append(words, p.peek().Literal.(string))
		p.advance()
	}
	for p.checkTokenType(COMMA) {
		p.advance()
		if p.checkTokenType(STRING) {
			words = append(words, p.peek().Literal.(string))
			p.advance()
		}
	}
	return words, err
}

func (p *Parser) parseComparison() (Expr, error) {
	var expr Expr
	var err error = nil
	var field string
	var OpToken Token
	var literal Literal

	if p.checkTokenType(IDENTIFIER) {
		field = p.peek().Lexeme
		p.advance()
	}
	if p.checkTokenType(OPERATOR) {
		OpToken = p.peek()
		p.advance()
	}
	if p.checkTokenType(STRING, IDENTIFIER, NUMBER) {
		val := p.peek().Literal
		literal = *NewLiteral(val)
		p.advance()
	}

	expr = NewComparison(field, OpToken, literal)
	return expr, err
}
