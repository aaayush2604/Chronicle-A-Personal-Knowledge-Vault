package parser

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	lexer "chronicle/internal/query/lexer"
	"fmt"
)

type Parser struct {
	Tokens  []*lexer.Token
	Current int
}

func NewParser(tokens []*lexer.Token) *Parser {
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

func (p *Parser) checkTokenType(tTypes ...lexer.TokenType) bool {
	for _, tType := range tTypes {
		if p.peek().TokenType == tType {
			return true
		}
	}
	return false
}

func (p *Parser) check(lexemes ...string) bool {
	if p.isAtEnd() {
		return false
	}
	for _, lexeme := range lexemes {
		return p.peek().Lexeme == lexeme
	}
	return false
}

func (p *Parser) advance() *lexer.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == lexer.EOF
}

func (p *Parser) peek() *lexer.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() *lexer.Token {
	return p.Tokens[p.Current-1]
}

// of no use yet, since only one expression tree is allowed
// func (p *Parser) synchronize() {
// 	p.advance()

// 	for !p.isAtEnd() {
// 		if p.check("recall") {
// 			return
// 		}

// 		if p.check("and") || p.check("or") {
// 			p.advance()
// 			return
// 		}

// 		p.advance()
// 	}
// }

func (p *Parser) Parse() (*Query, error) {
	q, err := p.parseQuery()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Query")
	}
	return q, nil
}

// Now we start adding a function to parse each of the grammar rules
func (p *Parser) parseQuery() (*Query, error) {
	var cmd CommandType
	c := p.peek()
	if c.TokenType != lexer.COMMAND {
		return nil, errorC.New(errorC.Syntax, "Syntax Error: Query must begin with a command")
	}
	switch c.Lexeme {
	case "recall":
		cmd = RecallCommand
	case "rem", "remember":
		cmd = RemCommand
	case "forget":
		cmd = ForgetCommand
	default:
		err := errorC.New(errorC.Syntax, "Syntax Error: Query must begin with a command")
		return nil, err
	}
	p.advance() //for consuming the cmd

	expr, ok := p.parseAll()
	if ok {
		if !p.isAtEnd() {
			return nil, errorC.New(errorC.Syntax, "There should be no input after ALL")
		}
		return &Query{
			Command: cmd,
			Expr:    expr,
		}, nil
	}

	queryNode := &Query{
		Command: cmd,
	}
	switch cmd {
	case RecallCommand:
		expr, err := p.parseExpression()

		if err != nil {
			err := errorC.Wrap(err, errorC.Syntax, "Error parsing Expression:")
			return nil, err
		}

		queryNode.Expr = expr
		return queryNode, nil

	case RemCommand:
		payload, err := p.parsePayload()

		if err != nil {
			err := errorC.Wrap(err, errorC.Syntax, "Error parsing Payload:")
			return nil, err
		}

		queryNode.Payload = payload
		return queryNode, nil
	case ForgetCommand:
		expr, err := p.parseExpression()

		if err != nil {
			err := errorC.Wrap(err, errorC.Syntax, "Error parsing Expression:")
			return nil, err
		}

		queryNode.Expr = expr
		return queryNode, nil
	}

	return nil, nil
}

func (p *Parser) parseAll() (Expr, bool) {
	if p.matchLexeme("all") {
		return &All{}, true
	}
	return nil, false
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
	if p.checkTokenType(lexer.LPAREN) {
		expr, err := p.parseGrouping()
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Grouping:")
		}
		return expr, nil
	} else if p.checkTokenType(lexer.IDENTIFIER) {
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
		} else if p.matchLexeme("tags") {
			expr, err := p.parseTags()
			if err != nil {
				return nil, errorC.Wrap(err, errorC.Syntax, "Error parsing Tags:")
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

	if !(p.checkTokenType(lexer.LBRACKET) && p.matchLexeme("[")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Contains StringList [ at col %d", p.peek().Position))
	}

	strings, err := p.parseStringList()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing StringList:")
	}
	expr := NewContains(strings)

	if !(p.checkTokenType(lexer.RBRACKET) && p.matchLexeme("]")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Contains StringList ] at col %d", p.peek().Position))
	}

	return expr, nil
}

func (p *Parser) parseTypeFilter() (Expr, error) {
	if !(p.checkTokenType(lexer.LBRACKET) && p.matchLexeme("[")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Type List [ at col %d", p.peek().Position))
	}

	strings, err := p.parseWordList()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing Type List:")
	}
	expr := NewTypeFilter(strings)

	if !(p.checkTokenType(lexer.RBRACKET) && p.matchLexeme("]")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Type List ] at col %d", p.peek().Position))
	}

	return expr, nil
}

func (p *Parser) parseTags() (Expr, error) {
	if !(p.checkTokenType(lexer.LBRACKET) && p.matchLexeme("[")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Tags List [ at col %d", p.peek().Position))
	}

	strings, err := p.parseWordList()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in parsing Tags List:")
	}
	expr := NewTags(strings)

	if !(p.checkTokenType(lexer.RBRACKET) && p.matchLexeme("]")) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Error parsing Tags List ] at col %d", p.peek().Position))
	}

	return expr, nil
}

func (p *Parser) parseStringList() ([]string, error) {
	if p.checkTokenType(lexer.RBRACKET) {
		return []string{}, nil
	}

	var words []string
	if !p.checkTokenType(lexer.STRING) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a String at col %d", p.peek().Position))
	}

	words = append(words, p.peek().Literal.(string))
	p.advance()

	for p.check(",") {
		p.advance()
		if !p.checkTokenType(lexer.STRING) {
			return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a String or ] at col %d", p.peek().Position))
		}
		words = append(words, p.peek().Literal.(string))
		p.advance()
	}
	return words, nil
}

func (p *Parser) parseWordList() ([]string, error) {
	if p.checkTokenType(lexer.RBRACKET) {
		return []string{}, nil
	}

	var words []string
	if !p.checkTokenType(lexer.IDENTIFIER) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Identifier at col %d", p.peek().Position))
	}

	words = append(words, p.peek().Lexeme)
	p.advance()

	for p.check(",") {
		p.advance()
		if !p.checkTokenType(lexer.IDENTIFIER) {
			return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Identifier or ] at col %d", p.peek().Position))
		}
		words = append(words, p.peek().Lexeme)
		p.advance()
	}
	return words, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	if !p.checkTokenType(lexer.IDENTIFIER) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Identifier at col: %d", p.peek().Position))
	}
	field := p.peek()
	p.advance()

	if !p.checkTokenType(lexer.OPERATOR) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected an Operator at col: %d", p.peek().Position))
	}
	opToken := p.peek()
	p.advance()

	if !p.checkTokenType(lexer.STRING, lexer.NUMBER) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected a Literal(Number or String) at col: %d", p.peek().Position))
	}
	literal := NewLiteral(p.peek())
	p.advance()

	return NewComparison(field, opToken, literal), nil
}

func (p *Parser) parsePayload() (Payload, error) {
	eType := entry.TypeNote
	if p.checkTokenType(lexer.ETYPE) {
		switch p.peek().Literal.(string) {
		case "n", "note":
			eType = entry.TypeNote
		case "l", "learning":
			eType = entry.TypeLearning
		case "q", "question":
			eType = entry.TypeQuestion
		case "i", "idea":
			eType = entry.TypeIdea
		case "imp", "important":
			eType = entry.TypeImportant
		}
		p.advance()
	}

	var tags []*lexer.Token
	for p.checkTokenType(lexer.TAG) {
		tags = append(tags, p.peek())
		p.advance()
	}

	if p.checkTokenType(lexer.COMMAND, lexer.ETYPE, lexer.TAG) {
		return nil, errorC.New(errorC.Syntax, fmt.Sprintf("Expected Entry Content at col %d", p.peek().Position))
	}

	var content []*lexer.Token
	for !p.checkTokenType(lexer.EOF) {
		content = append(content, p.peek())
		p.advance()
	}

	return &RemPayload{
		Type:    eType,
		Tags:    tags,
		Content: content,
	}, nil
}
