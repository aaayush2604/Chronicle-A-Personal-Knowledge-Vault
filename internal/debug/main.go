package main

import (
	"chronicle/internal/errorC"
	lexer "chronicle/internal/query/lexer"
	parser "chronicle/internal/query/parser"
	semantic "chronicle/internal/query/semantic"
	"fmt"
)

func main() {

	input := `recall date < "2022" AND time>"13" OR contains[] AND type["imp"]`
	scanner := lexer.NewScanner(input)

	tokens := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		fmt.Println(errorC.FormatError(err))
	}

	err = semantic.AnalyzeSemantics(q)
	if err != nil {
		errorString := errorC.FormatError(err)
		fmt.Println(errorString)
	}

	// for _, token := range tokens {
	// 	fmt.Printf("%s %v\n", token.Lexeme, token.TokenType)
	// }
}
