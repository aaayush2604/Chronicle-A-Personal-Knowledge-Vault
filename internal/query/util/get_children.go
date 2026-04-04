package util

import (
	"chronicle/internal/query/parser"
)

func GetChildren(expr parser.Expr) []parser.Expr {
	switch e := expr.(type) {
	case *parser.Logical:
		return []parser.Expr{e.Left, e.Right}
	case *parser.Grouping:
		return []parser.Expr{e.Expression}
	default:
		return []parser.Expr{}
	}
}
