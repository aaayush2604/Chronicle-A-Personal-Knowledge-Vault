package terminal

import (
	"chronicle/internal/query/lexer"
	"strings"
)

type CompletionContext string

const (
	RecallContext CompletionContext = "recall"
	RemContext    CompletionContext = "remember"
	TypeContext   CompletionContext = "type"
	TagContext    CompletionContext = "tag"
	NullContext   CompletionContext = "null"
)

type ChronicleCompleter struct {
}

func NewCompleter() *ChronicleCompleter {
	return &ChronicleCompleter{}
}

var isCompletedWord = map[string]struct{}{
	"recall":   {},
	"rem":      {},
	"remember": {},
	"[":        {},
	"]":        {},
	"(":        {},
	")":        {},
	"@":        {},
	"#":        {},
	"and":      {},
	"or":       {},
}

func isACompletedToken(s string) bool {
	_, ok := isCompletedWord[s]
	return ok
}

func determineCompletionContext(tokens []*lexer.Token) (CompletionContext, CompletionContext) {
	var commandContext CompletionContext = NullContext
	var context CompletionContext = NullContext

	for _, token := range tokens {
		if token.Lexeme == "recall" {
			commandContext = RecallContext
			context = NullContext
		} else if token.Lexeme == "rem" || token.Lexeme == "remember" {
			commandContext = RemContext
			context = NullContext
		} else if token.Lexeme == "type" || strings.HasPrefix(token.Lexeme, "@") {
			context = TypeContext
		} else if token.Lexeme == "tag" || strings.HasPrefix(token.Lexeme, "#") {
			context = TagContext
		} else if token.Lexeme == "and" || token.Lexeme == "or" || token.Lexeme == "(" {
			context = NullContext
		}
	}

	return commandContext, context
}

func getCompletionCandidates(commandContext CompletionContext, context CompletionContext) []string {
	if commandContext == RecallContext {
		switch context {
		case NullContext:
			return []string{"all", "contains[", "type[", "tags[", "time", "date", "len"}
		case TypeContext:
			return []string{"note", "learning", "question", "important", "idea"}
		case TagContext:
			return []string{}
		}
	} else if commandContext == RemContext {
		switch context {
		case NullContext:
			return []string{}
		case TypeContext:
			return []string{"@note", "@learning", "@question", "@important", "@idea"}
		case TagContext:
			return []string{}
		}
	}
	return []string{}
}

func filterCandidates(list []string, prefix string) []string {
	var res []string
	for _, l := range list {
		if strings.HasPrefix(l, prefix) {
			res = append(res, strings.TrimPrefix(l, prefix))
		}
	}

	return res
}

func toRunes(candidates []string) [][]rune {
	result := make([][]rune, 0, len(candidates))

	for _, candidate := range candidates {
		result = append(result, []rune(candidate))
	}

	return result
}

func (c *ChronicleCompleter) Do(line []rune, pos int) ([][]rune, int) {
	input := string(line[:pos])
	var prefix string

	if input == "" {
		return [][]rune{}, 0
	}

	scanner := lexer.NewScanner(input)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return [][]rune{}, 0
	}

	commandContext, context := determineCompletionContext(tokens)

	initialCandidates := getCompletionCandidates(commandContext, context)

	lastToken := tokens[len(tokens)-2]
	endsWithSpace := len(input) > 0 && input[len(input)-1] == ' '
	if endsWithSpace || isACompletedToken(lastToken.Lexeme) {
		prefix = ""
	} else {
		prefix = lastToken.Lexeme
	}

	candidates := filterCandidates(initialCandidates, prefix)

	return toRunes(candidates), len(prefix)
}
