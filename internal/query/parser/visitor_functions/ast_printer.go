package visitor

import (
	"chronicle/internal/query/parser"
	"fmt"
	"strings"
)

type ASTPrinter struct{}

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
		expr.Field,
		expr.Operator.Lexeme,
		expr.Value,
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
	return fmt.Sprintf("%v", expr.Val)
}

// func (p *ASTPrinter) PrintPreOrder(name string, exprs ...parser.Expr) string {
// 	var builder strings.Builder

// 	builder.WriteString("(")
// 	builder.WriteString(name)

// 	for _, expr := range exprs {
// 		builder.WriteString(" ")
// 		builder.WriteString(expr.Accept(p).(string))
// 	}
// 	builder.WriteString(")")
// 	return builder.String()
// }
