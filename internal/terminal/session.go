package terminal

import "fmt"

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner() {
	fmt.Println(bold + fgCyan + "Chronicle" + reset)
	fmt.Println(dim + "Personal Knowledge Log" + reset)
	fmt.Println()
}

func prompt() string {
	return fgCyan + "chronicle > " + reset
}
