package query

import (
	"unicode/utf8"
	"chronicle/internal/errorC"
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
	if s.current>=utf8.RuneCountInString(s.source)
}

func (s *Scanner) addToken(type TokenType){
	addToken(type,nil)
}

func (s *Scanner) addToken(type TokenType, literal any){
	text:=s.source[start,current]
	s.tokens=append(s.tokens,NewToken(type,text,literal,start))
}

func (s *Scanner) advance() rune {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() return false
	if s.source[s.current]==expected return false

	s.current++
	return true
}

func (s *Scanner) scanToken(){
	c:=s.advance()
	switch c {
	case '[':
		s.addToken(LBRACKET)
		break
	case ']':
		s.addToken(RBRACKET)
		break
	case ',':
		s.addToken(COMMA)
		break
	default:
		err:=errorC.New(NotFound, "Unexpected Charchter")
		fmt.Println(err.Error())
		break
	}
}

