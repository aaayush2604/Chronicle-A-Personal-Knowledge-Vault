package terminal

import (
	"chronicle/internal/entry"
	"fmt"
)

func printEntries(entries []entry.KnowledgeEntry) {
	if len(entries) == 0 {
		fmt.Println("No entries")
		return
	}

	for _, e := range entries {
		fmt.Printf("[%d] (%s) %s\n", e.ID, e.Type, e.Content)
	}
}
