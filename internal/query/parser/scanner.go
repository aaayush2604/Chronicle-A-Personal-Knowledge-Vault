package query

import (
	"chronicle/internal/errorC"
	"fmt"
	"strconv"
	"strings"
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
	if s.source[s.current] != expected {
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

func (s *Scanner) addString() {
	for !(s.peek() == '"') && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		err := errorC.New(errorC.Validation, fmt.Sprintf("Unterminated String at col %d", s.current))
		fmt.Println(err.Error())
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, value)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) addNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		e := errorC.New(errorC.Validation, fmt.Sprintf("Unable to parse Numeric Literal at col %d", s.current))
		fmt.Println(e.Error())
	}
	s.addTokenLiteral(NUMBER, value)
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) addIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := strings.ToUpper(s.source[s.start:s.current])
	tType, ok := keywords[text]
	if !ok {
		tType = IDENTIFIER
	}
	s.addToken(tType)
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
			s.addToken(OPERATOR)
		} else {
			s.addToken(OPERATOR)
		}
	case '>':
		if s.match('=') {
			s.addToken(OPERATOR)
		} else {
			s.addToken(OPERATOR)
		}
	case '=':
		s.addToken(OPERATOR)
	case '!':
		if s.match('=') {
			s.addToken(OPERATOR)
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
	case '"':
		s.addString()
	case ' ', '\r', '\t', '\n':
		return
	default:
		if s.isDigit(c) {
			s.addNumber()
		} else if s.isAlpha((c)) {
			s.addIdentifier()
		} else {
			err := errorC.New(errorC.NotFound, "Unexpected Character")
			fmt.Println(err.Error())
		}
	}
}
