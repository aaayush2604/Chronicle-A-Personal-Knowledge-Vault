package terminal

import (
	"chronicle/internal/entry"
	"fmt"
	"strconv"
	"strings"
)

func (r *REPL) handle(input string) bool {
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])

	var commandAliases = map[string]string{
		"n":   "note",
		"i":   "idea",
		"q":   "question",
		"l":   "learning",
		"imp": "important",
	}

	if full, ok := commandAliases[cmd]; ok {
		cmd = full
	}

	switch cmd {
	case "exit", "quit":
		return true

	case "help":
		printHelp(r.version)
		return false

	case "note", "idea", "question", "learning", "important":
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
		case "important":
			t = entry.TypeImportant
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

		if id, err := strconv.Atoi(parts[1]); err == nil {

			if deleted, ts := r.engine.CheckDelete(id); deleted {
				fmt.Printf(
					"Entry [%d] was deleted by user on %s\n",
					id,
					ts.Timestamp.Format("02 Jan 2006 15:04"),
				)
				return false
			}
		}
		var searchTerm string = strings.ToLower(parts[1])
		results := r.engine.Recall(searchTerm)
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
	case "version":
		fmt.Println("Chronicle v" + r.version)
		return false
	case "index":
		r.engine.PrintIndex()
		return false
	case "del":
		if len(parts) != 2 {
			fmt.Println("Usage: del <id>")
			return false
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid id")
			return false
		}

		if err := r.engine.Delete(id); err != nil {
			fmt.Println("Error:", err)
		}

		return false
	case "clear":
		clearScreen()
		return false
	default:
		fmt.Println("Unknown command. Type `help`.")
		return false
	}

}
