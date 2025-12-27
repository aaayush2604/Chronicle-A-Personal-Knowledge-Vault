package index

var stopWords = map[string]struct{}{
	"the": {}, "is": {}, "and": {}, "or": {}, "to": {},
	"of": {}, "in": {}, "on": {}, "for": {}, "with": {},
	"a": {}, "an": {}, "at": {}, "by": {}, "from": {},
}

func isStopWord(word string) bool {
	_, ok := stopWords[word]
	return ok
}
