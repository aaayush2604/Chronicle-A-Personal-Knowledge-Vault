package entry

import "strings"

type EntryType string

const (
	TypeNote      EntryType = "note"
	TypeIdea      EntryType = "idea"
	TypeQuestion  EntryType = "question"
	TypeLearning  EntryType = "learning"
	TypeImportant EntryType = "important"
)

var typeNames = map[string]EntryType{
	"note":      TypeNote,
	"n":         TypeNote,
	"idea":      TypeIdea,
	"i":         TypeIdea,
	"question":  TypeQuestion,
	"q":         TypeQuestion,
	"learning":  TypeLearning,
	"l":         TypeLearning,
	"important": TypeImportant,
	"imp":       TypeImportant,
}

func AllTypes() []EntryType {
	return []EntryType{TypeNote, TypeIdea, TypeQuestion, TypeLearning, TypeImportant}
}

func TypeFor(name string) (EntryType, bool) {
	t, ok := typeNames[strings.ToLower(strings.TrimSpace(name))]
	return t, ok
}

func Canonical(name string) EntryType {
	if t, ok := TypeFor(name); ok {
		return t
	}
	return EntryType(name)
}
