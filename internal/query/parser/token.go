package parser

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Position  int
}

func NewToken(tType TokenType, lex string, lit any, pos int) *Token {
	return &Token{
		TokenType: tType,
		Lexeme:    lex,
		Literal:   lit,
		Position:  pos,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.TokenType, t.Lexeme, t.Literal)
}
