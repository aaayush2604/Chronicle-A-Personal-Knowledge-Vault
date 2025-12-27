package index

import (
	"strings"
	"unicode"
)

func normalize(text string) string {
	text = strings.ToLower(text)

	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			b.WriteRune(r)
		}
	}

	return b.String()
}
