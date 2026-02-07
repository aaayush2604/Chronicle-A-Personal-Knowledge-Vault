package terminal

import (
	"fmt"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner(version string) {
	fmt.Println(bold + fgCyan + "Chronicle" + reset)
	fmt.Println(dim + "Personal Knowledge Vault - v" + version + reset)
	fmt.Println()
}

func prompt() {
	fmt.Print(fgCyan + "chronicle > " + reset)
}
