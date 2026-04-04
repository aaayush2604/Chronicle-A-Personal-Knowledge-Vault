package semantic

import (
	"fmt"

	"chronicle/internal/errorC"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/util"
)

type SemanticAnalyzer struct{}

func AnalyzeSemantics(q *parser.Query) error {
	analyzer := &SemanticAnalyzer{}
	res := q.Expr.Accept(analyzer)
	if err, ok := res.(error); ok && err != nil {
		return errorC.Wrap(err, errorC.Validation, "Error in Semantics: ")
	}
	return nil
}

func (s *SemanticAnalyzer) VisitLogicalExpression(expr *parser.Logical) any {
	Op := expr.Operator
	if Op.TokenType != lexer.LOGICAL {
		return errorC.New(errorC.Validation, fmt.Sprintf("Error in %v Semantics due to Incorrect Operator TokenType at col  %v", Op.Lexeme, Op.Position))
	}
	if Op.Lexeme != "and" && Op.Lexeme != "or" {
		return errorC.New(errorC.Validation, fmt.Sprintf("Invalid Logical Operator at col %v", Op.Position))
	}
	children := util.GetChildren(expr)
	for _, child := range children {
		res := child.Accept(s)
		if err, ok := res.(error); ok && err != nil {
			return errorC.Wrap(err, errorC.Validation, fmt.Sprintf("Error in '%v' Logical Semantics: ", Op.Lexeme))
		}
	}
	return nil
}

var validFieldTypes = map[string]struct{}{
	"date": {},
	"time": {},
	"len":  {},
}

func isValidFieldType(field string) bool {
	if _, ok := validFieldTypes[field]; !ok {
		return false
	}
	return true
}

var validComparisonOperators = map[string]struct{}{
	"<":  {},
	">":  {},
	"<=": {},
	">=": {},
	"=":  {},
}

func isValidComparisonOperator(Op string) bool {
	if _, ok := validComparisonOperators[Op]; !ok {
		return false
	}
	return true
}

func isValidComparison(field *lexer.Token, Value *lexer.Token) bool {
	if field.Lexeme == "date" || field.Lexeme == "time" {
		if Value.TokenType != lexer.STRING {
			return false
		}
	}
	if field.Lexeme == "len" {
		if Value.TokenType != lexer.NUMBER {
			return false
		}
	}
	return true
}

func (s *SemanticAnalyzer) VisitComparisonExpression(expr *parser.Comparison) any {
	field := expr.Field
	if field.TokenType != lexer.IDENTIFIER {
		return errorC.New(errorC.Validation, fmt.Sprintf("Error in Comparison Semantics due to Incorrect Field TokenType at col  %v", field.Position))
	}
	if !isValidFieldType(field.Lexeme) {
		return errorC.New(errorC.Validation, fmt.Sprintf("Field Not Allowed. Error in Comparison Semantics at col %v", field.Position))
	}
	Op := expr.Operator
	if Op.TokenType != lexer.OPERATOR {
		return errorC.New(errorC.Validation, fmt.Sprintf("Error in Comparison Semantics due to Incorrect Operator TokenType at col  %v", Op.Position))
	}
	if !isValidComparisonOperator(Op.Lexeme) {
		return errorC.New(errorC.Validation, fmt.Sprintf("Operator Not Allowed. Error in Comparison Semantics at col %v", Op.Position))
	}
	Value := expr.Value.Val
	if Value.TokenType != lexer.STRING && Value.TokenType != lexer.NUMBER {
		return errorC.New(errorC.Validation, fmt.Sprintf("Error in Comparison Semantics due to Incorrect Value TokenType at col  %v", Value.Position))
	}
	if !isValidComparison(field, Value) {
		return errorC.New(errorC.Validation, fmt.Sprintf("Invalid Comparison between %v and %v", field.Lexeme, Value.TokenType))
	}
	return nil
}

func (s *SemanticAnalyzer) VisitContainsExpression(expr *parser.Contains) any {
	words := expr.Words
	if len(words) == 0 {
		return errorC.New(errorC.Validation, "Empty Contains Not Allowed")
	}
	return nil
}

func (s *SemanticAnalyzer) VisitTypeFilterExpression(expr *parser.TypeFilter) any {
	words := expr.Words
	if len(words) == 0 {
		return errorC.New(errorC.Validation, "Empty Type Not Allowed")
	}
	return nil
}

func (s *SemanticAnalyzer) VisitGroupingExpression(expr *parser.Grouping) any {
	res := expr.Expression.Accept(s)
	if err, ok := res.(error); ok && err != nil {
		return errorC.Wrap(err, errorC.Validation, "Error in Grouping: ")
	}
	return nil
}

func (s *SemanticAnalyzer) VisitLiteralExpression(expr *parser.Literal) any {
	return nil
}
