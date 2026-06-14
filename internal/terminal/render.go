package terminal

import (
	"chronicle/internal/entry"
	"fmt"
	"strings"
)

func printEntries(entries []entry.KnowledgeEntry) {
	if len(entries) == 0 {
		fmt.Println(dim + "No entries" + reset)
		return
	}

	count := 0
	for _, e := range entries {
		// fmt.Printf("[%d] (%s) %s\n", e.ID, e.Type, e.Content)

		var sb strings.Builder
		for i, t := range e.Tags {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, "#%s", t)
		}

		fmt.Printf(
			"%s{%d}%s %s%s%s (%s)  [%s] %s \n\n",
			fgGray, e.ID, reset,
			colorForType(e.Type), e.Content, reset,
			e.Type, e.Timestamp.Format("3:04 PM 02-01-2006"), sb.String(),
		)

		count++
		if count%pageSize == 0 {
			fmt.Print(dim + "Press Enter to continue..." + reset)
			pause()
		}
	}
}

func colorForType(t entry.EntryType) string {
	switch t {
	case entry.TypeIdea:
		return fgYellow
	case entry.TypeQuestion:
		return fgBlue
	case entry.TypeLearning:
		return fgGreen
	case entry.TypeImportant:
		return fgRed
	default:
		return fgWhite
	}
}
