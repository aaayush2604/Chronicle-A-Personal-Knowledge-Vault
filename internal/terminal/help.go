package terminal

import "fmt"

func printHelp(version string) {
	fmt.Println(dim + "Chronicle v" + version + reset)
	fmt.Println()

	// ---------------- WRITING ----------------
	fmt.Println(bold + "Writing" + reset)
	fmt.Println("  n/note <text>")
	fmt.Println("  i/idea <text>")
	fmt.Println("  q/question <text>")
	fmt.Println("  l/learning <text>")
	fmt.Println("  imp <text>")
	fmt.Println()

	// ---------------- RECALL ----------------
	fmt.Println(bold + "Recall & Query Language" + reset)
	fmt.Println("  recall <query>")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    recall <predicate> AND/OR <predicate> AND/OR <predicate> .....")

	// ---------------- QUERY SYNTAX ----------------
	fmt.Println(bold + "Query Syntax" + reset)
	fmt.Println("  Fields:")
	fmt.Println("    date, time, contains, type")
	fmt.Println()

	fmt.Println("  Comparison Operators:")
	fmt.Println("    =   equal")
	fmt.Println("    !=  not equal")
	fmt.Println("    >   greater than")
	fmt.Println("    <   less than")
	fmt.Println("    >=  greater or equal")
	fmt.Println("    <=  less or equal")
	fmt.Println()

	fmt.Println("  Logical Operators:")
	fmt.Println("    AND   OR")
	fmt.Println()

	fmt.Println("  Examples:")
	fmt.Println("    recall type[\"imp\",\"note\"] (Satisfied for any one of them)")
	fmt.Println("    recall time > \"7 PM\" AND type[\"note\"]")
	fmt.Println("    recall date >= \"2026-04-01\"")
	fmt.Println("    recall contains [\"string\",\"string\"] (has to satisfy all of the strings)")

	fmt.Println()

	// ---------------- TIME & DATE ----------------
	fmt.Println(bold + "Time & Date Formats" + reset)
	fmt.Println("  Time:")
	fmt.Println("    \"7 PM\", \"07:30 PM\", \"19:30\"")
	fmt.Println()
	fmt.Println("  Date:")
	fmt.Println("    \"2026-04-10\", \"10/04/2026\"")
	fmt.Println()

	// ---------------- QUICK FILTERS ----------------
	fmt.Println(bold + "Quick Commands" + reset)
	fmt.Println("  today")
	fmt.Println("  this week")
	fmt.Println("  month")
	fmt.Println("  year")
	fmt.Println("  summary")
	fmt.Println()

	// ---------------- DELETION ----------------
	fmt.Println(bold + "Deletion" + reset)
	fmt.Println("  delete <id>")
	fmt.Println()

	// ---------------- SESSION ----------------
	fmt.Println(bold + "Session" + reset)
	fmt.Println("  help")
	fmt.Println("  exit | quit")
}
