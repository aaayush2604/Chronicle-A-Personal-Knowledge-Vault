package lexer

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

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return []*Token{}, err
		}
	}

	endToken := NewToken(EOF, "", nil, 0)
	s.tokens = append(s.tokens, endToken)
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tType TokenType) {
	s.addTokenLiteral(tType, nil)
}

func (s *Scanner) addTokenLiteral(tType TokenType, literal any) {
	text := s.source[s.start:s.current]
	text = strings.ToLower(text)
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

func (s *Scanner) addString() error {
	for !(s.peek() == '"') && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		err := errorC.New(errorC.Validation, fmt.Sprintf("Unterminated String at col %d", s.current))
		return err
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, value)
	return nil
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) addNumber() error {
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
		return e
	}
	s.addTokenLiteral(NUMBER, value)
	return nil
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) addIdentifier() error {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := strings.ToUpper(s.source[s.start:s.current])
	tType, ok := keywords[text]
	if !ok {
		tType = IDENTIFIER
	}
	s.addToken(tType)
	return nil
}

func (s *Scanner) addTag() error {
	if s.peek() == '#' {
		return errorC.New(errorC.Validation, fmt.Sprintf("Unexpected '#' at col %d", s.current))
	}
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	tagText := strings.ToLower(s.source[s.start+1 : s.current])
	s.addTokenLiteral(TAG, tagText)
	return nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '[':
		s.addToken(LBRACKET)
	case ']':
		s.addToken(RBRACKET)
	case '(':
		s.addToken(LPAREN)
	case ')':
		s.addToken(RPAREN)
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
				return err
			}

			s.advance()
			s.advance()
		} else {
			err := errorC.New(errorC.NotFound, "Unterminated comment")
			return err
		}
	case '"':
		err := s.addString()
		if err != nil {
			return errorC.Wrap(err, errorC.Syntax, "Error parsing Query: ")
		}
	case '#':
		err := s.addTag()
		if err != nil {
			return errorC.Wrap(err, errorC.Syntax, "Error parsing Query: ")
		}
	case ' ', '\r', '\t', '\n':
		return nil
	default:
		if s.isDigit(c) {
			err := s.addNumber()
			if err != nil {
				return errorC.Wrap(err, errorC.Syntax, "Error parsing Query: ")
			}
		} else if s.isAlpha((c)) {
			err := s.addIdentifier()
			if err != nil {
				return errorC.Wrap(err, errorC.Syntax, "Error parsing Query: ")
			}
		} else {
			err := errorC.New(errorC.NotFound, "Unexpected Character")
			return err
		}
	}
	return nil
}
