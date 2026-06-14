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

func TestNewRemType(t *testing.T) {
	n := NewRemember(nil)

	if n.GetType() != RemType {
		t.Fatalf("expected rem type")
	}
}
