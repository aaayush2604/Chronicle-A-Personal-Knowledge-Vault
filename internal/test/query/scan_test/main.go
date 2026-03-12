package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	query "chronicle/internal/query/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter query (type 'exit' to quit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			return
		}

		scanner := query.NewScanner(input)
		tokens := scanner.ScanTokens()

		fmt.Println("\nTokens:")
		for _, token := range tokens {
			fmt.Println(token)
		}

		fmt.Println()
	}
}
