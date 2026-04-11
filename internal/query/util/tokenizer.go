package util

import (
	"strings"
	"unicode"
)

var stopWords = map[string]struct{}{
	"the": {}, "is": {}, "and": {}, "or": {}, "to": {},
	"of": {}, "in": {}, "on": {}, "for": {}, "with": {},
	"a": {}, "an": {}, "at": {}, "by": {}, "from": {},
}

func isStopWord(word string) bool {
	_, ok := stopWords[word]
	return ok
}

func Tokenize(text string) []string {
	text = normalize(text)

	raw := strings.Fields(text)
	tokens := make([]string, 0, len(raw))

	for _, w := range raw {
		if len(w) < 2 {
			continue
		}
		if isStopWord(w) {
			continue
		}
		tokens = append(tokens, w)
	}

	return tokens
}

func normalize(text string) string {
	text = strings.ToLower(text)

	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}

	return b.String()
}
