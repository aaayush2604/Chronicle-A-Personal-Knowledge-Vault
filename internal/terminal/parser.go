package terminal

import (
	"chronicle/internal/entry"
	"fmt"
	"strings"
)

func (r *REPL) handle(input string) bool {
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "exit", "quit":
		return true

	case "help":
		printHelp()
		return false

	case "note", "idea", "question", "learning":
		if len(parts) < 2 {
			fmt.Printf("Usage: %s <text>\n", cmd)
			return false
		}

		content := strings.Join(parts[1:], " ")

		var t entry.EntryType
		switch cmd {
		case "note":
			t = entry.TypeNote
		case "idea":
			t = entry.TypeIdea
		case "question":
			t = entry.TypeQuestion
		case "learning":
			t = entry.TypeLearning
		}

		e, err := r.engine.AddNote(content, t)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
		fmt.Printf("Saved [%d] (%s)\n", e.ID, t)
		return false

	case "add":
		if len(parts) < 2 {
			fmt.Printf("Usage:add <text> [default-NOTE]\n")
			return false
		}

		content := strings.Join(parts[1:], " ")

		var t entry.EntryType = entry.TypeNote

		e, err := r.engine.AddNote(content, t)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
		fmt.Printf("Saved [%d] (%s)\n", e.ID, t)
		return false

	case "recall":
		if len(parts) < 2 {
			fmt.Println("Usage: recall <word>")
			return false
		}

		results := r.engine.Recall(parts[1])
		if len(results) == 0 {
			fmt.Println("No results")
			return false
		}

		printEntries(results)

		return false

	case "today":
		results := r.engine.Today()
		printEntries(results)
		return false

	case "this", "week":
		results := r.engine.ThisWeek()
		printEntries(results)
		return false

	case "summary":
		results := r.engine.ThisWeek()
		summary := r.engine.SummaryByType(results)

		printSummary(summary)
		return false

	default:
		fmt.Println("Unknown command. Type `help`.")
		return false
	}

}
