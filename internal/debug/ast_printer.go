package main

import (
	parser "chronicle/internal/query/parser"
	"fmt"
	"strings"
)

type ASTPrinter struct{}

func PrintAST(q *parser.Query) {
	printer := &ASTPrinter{}
	fmt.Println(q.Command)
	fmt.Println(q.Expr.Accept(printer).(string))
}

func (p *ASTPrinter) VisitLogicalExpression(expr *parser.Logical) any {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(expr.Left.Accept(p).(string))
	builder.WriteString(expr.Operator.Lexeme)
	builder.WriteString(expr.Right.Accept(p).(string))
	builder.WriteString(")")

	return builder.String()
}

func (p *ASTPrinter) VisitComparisonExpression(expr *parser.Comparison) any {
	return fmt.Sprintf("(%s %s %v)",
		expr.Field.Lexeme,
		expr.Operator.Lexeme,
		expr.Value.Val.Lexeme,
	)
}

func (p *ASTPrinter) VisitGroupingExpression(expr *parser.Grouping) any {
	return fmt.Sprintf("(group %s)", expr.Expression.Accept(p).(string))
}

func (p *ASTPrinter) VisitContainsExpression(expr *parser.Contains) any {
	var builder strings.Builder

	builder.WriteString("( contains [")
	for i, str := range expr.Words {
		builder.WriteString(str)
		if i < len(expr.Words)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("])")
	return builder.String()
}

func (p *ASTPrinter) VisitTypeFilterExpression(expr *parser.TypeFilter) any {
	var builder strings.Builder

	builder.WriteString("( contains [")
	for i, str := range expr.Words {
		builder.WriteString(str)
		if i < len(expr.Words)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("])")
	return builder.String()
}

func (p *ASTPrinter) VisitLiteralExpression(expr *parser.Literal) any {
	return fmt.Sprintf("%v", expr.Val.Lexeme)
}

func (p *ASTPrinter) VisitAllExpression(expr *parser.All) any {
	return "ALL"
}
