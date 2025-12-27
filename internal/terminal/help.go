package terminal

import "fmt"

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  note <text>       Add a note")
	fmt.Println("  idea <text>       Add an idea")
	fmt.Println("  question <text>   Add a question")
	fmt.Println("  learning <text>   Add a learning")
	fmt.Println("  recall <word>     Recall related entries")
	fmt.Println("  help              Show this help")
	fmt.Println("  exit | quit       Exit Chronicle")
	fmt.Println("  today             Review today’s entries")
	fmt.Println("  this week         Review entries from last 7 days")
	fmt.Println("  summary           Summary of last week by type")
}
