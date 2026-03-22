package main

import (
	lexer "chronicle/internal/query/lexer"
	parser "chronicle/internal/query/parser"
	"fmt"
)

func main() {

	input := `recall date < "2022" AND time>13 OR contains["ai","aayush"] AND type["learing"]`
	scanner := lexer.NewScanner(input)

	tokens := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		fmt.Println(err.Error())
	}
	PrintAST(q)
}
