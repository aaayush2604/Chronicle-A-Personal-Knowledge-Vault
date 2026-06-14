package terminal

import "fmt"

func printHelp(version string) {
	fmt.Println(dim + "Chronicle v" + version + reset)
	fmt.Println()

	// ---------------- WRITING ----------------
	fmt.Println(bold + "Remembering" + reset)
	fmt.Println("  rem [@type] [#tags...] <text>")
	fmt.Println("  remember [@type] [#tags...] <text>")
	fmt.Println()
	fmt.Println("  Types:")
	fmt.Println("    note (n)")
	fmt.Println("    idea (i)")
	fmt.Println("    question (q)")
	fmt.Println("    learning (l)")
	fmt.Println("    important (imp)")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    rem \"Read chapter 3\"")
	fmt.Println("    rem @learning #math #analysis \"Uniform convergence theorem\"")
	fmt.Println("    remember @idea #project \"Use an append-only log\"")
	fmt.Println()

	// ---------------- RECALL ----------------
	fmt.Println(bold + "Recall" + reset)
	fmt.Println("  recall <query>")
	fmt.Println()

	// ---------------- QUERY SYNTAX ----------------
	fmt.Println(bold + "Predicates" + reset)
	fmt.Println("  all")
	fmt.Println("  contains[...]")
	fmt.Println("  type[...]")
	fmt.Println("  tags[...]")
	fmt.Println("  date")
	fmt.Println("  time")
	fmt.Println("  len")
	fmt.Println()

	fmt.Println(bold + "Comparison Operators" + reset)
	fmt.Println("  =   equal")
	fmt.Println("  !=  not equal")
	fmt.Println("  >   greater than")
	fmt.Println("  <   less than")
	fmt.Println("  >=  greater or equal")
	fmt.Println("  <=  less or equal")
	fmt.Println()

	fmt.Println(bold + "Logical Operators" + reset)
	fmt.Println("  AND")
	fmt.Println("  OR")
	fmt.Println()

	fmt.Println(bold + "Examples" + reset)
	fmt.Println("  recall all")
	fmt.Println("  recall contains[\"golang\"]")
	fmt.Println("  recall contains[\"golang\"] AND type[note]")
	fmt.Println("  recall tags[go]")
	fmt.Println("  recall tags[go,database]")
	fmt.Println("  recall type[note,idea]")
	fmt.Println("  recall len > 100")
	fmt.Println("  recall time > \"7 PM\"")
	fmt.Println("  recall date >= \"2026-04-01\"")
	fmt.Println()

	fmt.Println("  Lists use OR semantics:")
	fmt.Println("    type[note,idea]")
	fmt.Println("    tags[go,database]")
	fmt.Println("    contains[\"go\",\"database\"]")
	fmt.Println()

	fmt.Println("  Require multiple conditions using AND:")
	fmt.Println("    tags[go] AND tags[database]")
	fmt.Println("    contains[\"go\"] AND contains[\"database\"]")
	fmt.Println()

	// ---------------- TIME & DATE ----------------
	fmt.Println(bold + "Time & Date Formats" + reset)
	fmt.Println("  Time:")
	fmt.Println("    \"7 PM\"")
	fmt.Println("    \"07:30 PM\"")
	fmt.Println("    \"19:30\"")
	fmt.Println()

	fmt.Println("  Date:")
	fmt.Println("    \"2026-04-10\"")
	fmt.Println("    \"10/04/2026\"")
	fmt.Println()

	// ---------------- QUICK COMMANDS ----------------
	fmt.Println(bold + "Quick Commands" + reset)
	fmt.Println("  today")
	fmt.Println("  this week")
	fmt.Println("  month")
	fmt.Println("  year")
	fmt.Println("  summary")
	fmt.Println()

	// ---------------- DELETION ----------------
	fmt.Println(bold + "Deletion" + reset)
	fmt.Println("  del <id>")
	fmt.Println()

	// ---------------- SESSION ----------------
	fmt.Println(bold + "Session" + reset)
	fmt.Println("  help")
	fmt.Println("  exit")
	fmt.Println("  quit")

}
