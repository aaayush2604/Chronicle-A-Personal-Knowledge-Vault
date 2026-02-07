package terminal

import "fmt"

func printHelp(version string) {
	fmt.Println(dim + "Chronicle v" + version + reset)
	fmt.Println()

	fmt.Println(bold + "Writing" + reset)
	fmt.Println("  n/note <text>")
	fmt.Println("  i/idea <text>")
	fmt.Println("  q/question <text>")
	fmt.Println("  l/learning <text>")
	fmt.Println("  imp <text>")
	fmt.Println()

	fmt.Println(bold + "Recall & Review" + reset)
	fmt.Println("  recall <word>")
	fmt.Println("  today")
	fmt.Println("  this week")
	fmt.Println("  month")
	fmt.Println("  year")
	fmt.Println("  summary")
	fmt.Println()

	fmt.Println(bold + "Deletion" + reset)
	fmt.Println("  delete <id>")
	fmt.Println()

	fmt.Println(bold + "Session" + reset)
	fmt.Println("  help")
	fmt.Println("  exit | quit")
}
