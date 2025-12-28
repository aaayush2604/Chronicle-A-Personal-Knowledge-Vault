package index

import (
	"strings"
	"unicode"
)

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
