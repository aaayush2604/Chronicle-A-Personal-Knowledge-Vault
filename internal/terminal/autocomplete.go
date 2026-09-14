package terminal

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
	"strings"
)

type CommandContext string

const (
	NoCommand     CommandContext = ""
	RecallContext CommandContext = "recall"
	RemContext    CommandContext = "remember"
	ForgetContext CommandContext = "forget"
	ReviseContext CommandContext = "revise"
)

type ChronicleCompleter struct {
}

func NewCompleter() *ChronicleCompleter {
	return &ChronicleCompleter{}
}

var commandNames = []string{
	"recall", "remember", "rem", "revise", "forget",
	"today", "week", "month", "year", "summary",
	"help", "index", "version", "clear", "exit", "quit",
}

var predicates = []string{
	"all", "contains[", "type[", "tags[", "id[", "time", "date", "len",
}

var connectors = []string{"AND", "OR"}

func typeCandidates(prefix string) []string {
	types := entry.AllTypes()
	out := make([]string, 0, len(types))
	for _, t := range types {
		out = append(out, prefix+string(t))
	}
	return out
}

type clause struct {
	command    CommandContext
	afterWhere bool
	listOwner  string
	hasType    bool
	inContent  bool
	prev       *lexer.Token
}

func isBoundary(r rune) bool {
	switch r {
	case ' ', '\t', '[', ']', '(', ')', ',', '"':
		return true
	}
	return false
}

func currentWord(line []rune) string {
	i := len(line)
	for i > 0 && !isBoundary(line[i-1]) {
		i--
	}
	return string(line[i:])
}

func readClause(tokens []*lexer.Token, word string) clause {
	if n := len(tokens); n > 0 && tokens[n-1].TokenType == lexer.EOF {
		tokens = tokens[:n-1]
	}

	if word != "" && len(tokens) > 0 {
		tokens = tokens[:len(tokens)-1]
	}

	c := clause{command: NoCommand}
	for _, t := range tokens {
		switch t.TokenType {
		case lexer.ETYPE:
			c.hasType = true
		case lexer.TEXT:
			c.inContent = true
		}
	}

	var lastIdentifier string

	for _, t := range tokens {
		switch t.TokenType {
		case lexer.LBRACKET:
			c.listOwner = lastIdentifier
		case lexer.RBRACKET:
			c.listOwner = ""
		case lexer.IDENTIFIER:
			lastIdentifier = t.Lexeme
		}

		if t.TokenType != lexer.COMMAND {
			continue
		}

		switch t.Lexeme {
		case "recall":
			c.command = RecallContext
		case "rem", "remember":
			c.command = RemContext
		case "forget":
			c.command = ForgetContext
		case "revise":
			c.command = ReviseContext
		case "where":
			c.afterWhere = true
		}
	}

	if n := len(tokens); n > 0 {
		c.prev = tokens[n-1]
	}
	return c
}

func endsPredicate(t *lexer.Token) bool {
	if t == nil {
		return false
	}
	switch t.TokenType {
	case lexer.RBRACKET, lexer.RPAREN, lexer.NUMBER, lexer.STRING, lexer.ALL:
		return true
	}
	return false
}

func predicateCandidates(c clause) []string {
	if endsPredicate(c.prev) {
		return connectors
	}
	return predicates
}

func candidatesFor(c clause) []string {
	if c.listOwner != "" {
		switch c.listOwner {
		case "type":
			return typeCandidates("")
		}
		return nil
	}

	switch c.command {
	case NoCommand:
		return commandNames

	case RecallContext, ForgetContext:
		return predicateCandidates(c)

	case RemContext:
		if c.inContent || c.hasType {
			return nil
		}
		if c.prev != nil && (c.prev.TokenType == lexer.COMMAND || c.prev.TokenType == lexer.TAG) {
			return typeCandidates("@")
		}
		return nil

	case ReviseContext:
		if c.afterWhere {
			return predicateCandidates(c)
		}
		if c.prev != nil && c.prev.TokenType == lexer.COMMAND {
			return typeCandidates("@")
		}
		return []string{"where"}
	}

	return nil
}

func filterCandidates(list []string, word string) []string {
	var res []string
	for _, l := range list {
		if l != word && strings.HasPrefix(l, word) {
			res = append(res, strings.TrimPrefix(l, word))
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
	if pos < 0 {
		pos = 0
	}
	if pos > len(line) {
		pos = len(line)
	}
	head := line[:pos]

	tokens, err := lexer.NewScanner(string(head)).ScanTokens()
	if err != nil {
		return [][]rune{}, 0
	}

	word := currentWord(head)
	candidates := filterCandidates(candidatesFor(readClause(tokens, word)), word)

	return toRunes(candidates), len([]rune(word))
}
