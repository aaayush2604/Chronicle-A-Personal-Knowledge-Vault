package index

import "strings"

func tokenize(text string) []string {
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
