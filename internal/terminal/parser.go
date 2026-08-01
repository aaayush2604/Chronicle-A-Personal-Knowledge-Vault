package terminal

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	"fmt"
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

	case "rem", "remember":
		if len(parts) < 2 {
			fmt.Printf("Usage: note <entry type (optional)> <#tags (optional)> <text>\n")
			return false
		}

		results, err := r.engine.Query(input)
		if err != nil {
			fmt.Println(errorC.FormatError(err))
			return false
		}

		fmt.Printf("Saved [%d] (%s)\n", results[0].ID, results[0].Type)

		return false

	case "recall":
		if len(parts) < 2 {
			fmt.Println("Usage: recall <predicate>")
			return false
		}

		results, err := r.engine.Query(input)
		if err != nil {
			fmt.Println(errorC.FormatError(err))
		}
		printEntries(results)

		return false
	case "forget":
		if len(parts) < 2 {
			fmt.Println("Usage: forget <predicate>")
			return false
		}

		deletionEntries, err := r.engine.Query(input)
		if err != nil {
			fmt.Println(errorC.FormatError(err))
		}

		delete, toBeDeleted := r.confirmDeletion(deletionEntries)

		if delete {
			err = r.engine.ProcessDeletion(toBeDeleted, deletionEntries)
			if err != nil {
				fmt.Println(errorC.FormatError(err))
			}

			var deletedCount int
			if len(toBeDeleted) > 0 {
				deletedCount = len(toBeDeleted)
			} else {
				deletedCount = len(deletionEntries)
			}
			fmt.Printf("[%d] entries deleted\n", deletedCount)
		} else {
			fmt.Printf("Forget Cancelled\n")
		}

		return false
	case "today":
		results := r.engine.Today()
		printEntries(results)
		return false

	case "month":
		results := r.engine.ThisMonth()
		printEntries(results)
		return false

	case "year":
		results := r.engine.ThisYear()
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
	case "clear":
		clearScreen()
		return false
	default:
		fmt.Println("Unknown command. Type `help`.")
		return false
	}

}

func (r *REPL) confirmDeletion(entries []entry.KnowledgeEntry) (bool, []int) {
	oldPrompt := r.prompt
	var response bool
	var list []int

	if len(entries) == 0 {
		fmt.Println("No Matching Entries Found!")
		return false, []int{}
	}

	printEntries(entries)

	r.terminal.SetPrompt("Delete? [Y/N/ id1,id2,id3...] > ")
	defer r.terminal.SetPrompt(oldPrompt)
	for {
		answer, err := r.terminal.ReadLine()
		if err != nil {
			return false, []int{}
		}
		if stop, r, l := acceptDeletionInput(answer); stop {
			response = r
			list = l
			break
		}
	}

	return response, list
}

func acceptDeletionInput(answer string) (bool, bool, []int) {
	if strings.ToLower(strings.TrimSpace(answer)) == "y" || strings.ToLower(strings.TrimSpace(answer)) == "yes" {
		return true, true, []int{}
	}

	if strings.ToLower(strings.TrimSpace(answer)) == "n" || strings.ToLower(strings.TrimSpace(answer)) == "no" {
		return true, false, []int{}
	}

	if ok, list := isIntList(answer); ok {
		return true, true, list
	}

	return false, false, []int{}
}
