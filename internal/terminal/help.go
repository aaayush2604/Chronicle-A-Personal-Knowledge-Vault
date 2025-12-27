package terminal

import "fmt"

func printHelp() {
	fmt.Println(bold + "Writing" + reset)
	fmt.Println("  note <text>")
	fmt.Println("  idea <text>")
	fmt.Println("  question <text>")
	fmt.Println("  learning <text>")
	fmt.Println()

	fmt.Println(bold + "Recall & Review" + reset)
	fmt.Println("  recall <word>")
	fmt.Println("  today")
	fmt.Println("  this week")
	fmt.Println("  summary")
	fmt.Println()

	fmt.Println(bold + "Session" + reset)
	fmt.Println("  help")
	fmt.Println("  exit | quit")
}
