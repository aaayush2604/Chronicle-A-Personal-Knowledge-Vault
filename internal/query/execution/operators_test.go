package execution

import (
	"chronicle/internal/query/parser"
	"testing"
)

func TestNewRecallType(t *testing.T) {
	r := NewRecall(&parser.All{})

	if r.GetType() != RecallType {
		t.Fatalf("expected recall type")
	}
}

func TestNewNoteType(t *testing.T) {
	n := NewNote(nil)

	if n.GetType() != NoteType {
		t.Fatalf("expected note type")
	}
}
