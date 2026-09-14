package lexer

import (
	"sort"
	"strings"
)

var keywords = map[string]TokenType{
	"AND":      LOGICAL,
	"OR":       LOGICAL,
	"RECALL":   COMMAND,
	"REM":      COMMAND,
	"REMEMBER": COMMAND,
	"FORGET":   COMMAND,
	"REVISE":   COMMAND,
	"WHERE":    COMMAND,
	"ALL":      ALL,
}

func ReservedWords() []string {
	words := make([]string, 0, len(keywords))
	for word := range keywords {
		words = append(words, strings.ToLower(word))
	}
	sort.Strings(words)
	return words
}

func IsReserved(t *Token) bool {
	switch t.TokenType {
	case ALL, LOGICAL, COMMAND:
		return true
	}
	return false
}
