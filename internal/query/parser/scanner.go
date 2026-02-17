package query

import "unicode/utf8"

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
