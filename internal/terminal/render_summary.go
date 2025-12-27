package terminal

import (
	"fmt"

	"chronicle/internal/entry"
)

func printSummary(summary map[entry.EntryType]int) {
	if len(summary) == 0 {
		fmt.Println(dim + "No entries in this period." + reset)
		return
	}

	fmt.Println(bold + "Summary · last 7 days" + reset)
	fmt.Println()

	// stable order for readability
	order := []entry.EntryType{
		entry.TypeIdea,
		entry.TypeQuestion,
		entry.TypeLearning,
		entry.TypeNote,
	}

	for _, t := range order {
		count, ok := summary[t]
		if !ok || count == 0 {
			continue
		}

		fmt.Printf(
			"%s%-10s%s %d\n",
			fgGray,
			t,
			reset,
			count,
		)
	}
}
