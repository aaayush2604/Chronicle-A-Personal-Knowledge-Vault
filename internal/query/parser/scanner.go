package query

import (
	"chronicle/internal/errorC"
	"fmt"
)

type Scanner struct {
	source  string
	tokens  []*Token
	start   int
	current int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		source: src,
	}
}

func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	endToken := NewToken(EOF, "", nil, 0)
	s.tokens = append(s.tokens, endToken)
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tType TokenType) {
	s.addTokenLiteral(tType, nil)
}

func (s *Scanner) addTokenLiteral(tType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tType, text, literal, s.start))
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] == expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '[':
		s.addToken(LBRACKET)
	case ']':
		s.addToken(RBRACKET)
	case ',':
		s.addToken(COMMA)
	case '<':
		if s.match('=') {
			s.addTokenLiteral(OPERATOR, "<=")
		} else {
			s.addTokenLiteral(OPERATOR, "<")
		}
	case '>':
		if s.match('=') {
			s.addTokenLiteral(OPERATOR, ">=")
		} else {
			s.addTokenLiteral(OPERATOR, ">")
		}
	case '!':
		if s.match('!') {
			s.addTokenLiteral(OPERATOR, "!=")
		} else {
			err := errorC.New(errorC.NotFound, "Unexpected Character")
			fmt.Println(err.Error())
		}
	case '/':
		if s.match('*') {
			for !(s.peek() == '*' && s.peekNext() == '/') && !s.isAtEnd() {
				s.advance()
			}

			if s.isAtEnd() {
				err := errorC.New(errorC.NotFound, "Unterminated comment")
				fmt.Println(err.Error())
				return
			}

			s.advance()
			s.advance()
		} else {
			err := errorC.New(errorC.NotFound, "Unterminated comment")
			fmt.Println(err.Error())
		}
	default:
		err := errorC.New(errorC.NotFound, "Unexpected Character")
		fmt.Println(err.Error())
	}
}
