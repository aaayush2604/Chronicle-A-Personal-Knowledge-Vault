package execution

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/util"
	"slices"
	"strings"
	"time"
)

type Record map[string]any

type AstEvaluator struct {
	record Record
}

func EvaluateRecall(expr parser.Expr, record Record) bool {
	e := &AstEvaluator{record: record}
	return expr.Accept(e).(bool)
}

func (e *AstEvaluator) VisitLogicalExpression(expr *parser.Logical) any {
	left := expr.Left.Accept(e).(bool)

	op := strings.ToLower(expr.Operator.Lexeme)

	if op == "or" {
		if left {
			return true
		}
		return expr.Right.Accept(e).(bool)
	} else if op == "and" {
		if !left {
			return false
		}
		return expr.Right.Accept(e).(bool)
	}

	return false
}

func (e *AstEvaluator) VisitComparisonExpression(expr *parser.Comparison) any {
	field := expr.Field.Lexeme
	op := expr.Operator.Lexeme

	if field == "len" {
		recordVal, ok := e.record["len"].(int)
		if !ok {
			return false
		}
		queryVal, ok := toFloat(expr.Value.Accept(e))
		if !ok {
			return false
		}
		return compareNumber(op, float64(recordVal), queryVal)
	}

	t, ok := e.record["date"].(time.Time)
	if !ok {
		return false
	}

	value, ok := expr.Value.Accept(e).(string)
	if !ok {
		return false
	}

	var start, end, recordVal int
	var err error

	switch field {
	case "date":
		start, end, err = util.ParseDate(value)
		recordVal = util.RecordDate(t)
	case "time":
		start, end, err = util.ParseTime(value)
		recordVal = util.RecordTime(t)
	default:
		return false
	}

	if err != nil {
		return false
	}

	return compareRange(op, recordVal, start, end)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}

func compareNumber(op string, recordVal, queryVal float64) bool {
	switch op {
	case "=":
		return recordVal == queryVal
	case "!=":
		return recordVal != queryVal
	case ">":
		return recordVal > queryVal
	case ">=":
		return recordVal >= queryVal
	case "<":
		return recordVal < queryVal
	case "<=":
		return recordVal <= queryVal
	}
	return false
}

func compareRange(op string, recordVal, start, end int) bool {
	switch op {
	case "=":
		return recordVal >= start && recordVal <= end
	case "!=":
		return recordVal < start || recordVal > end
	case ">":
		return recordVal > start
	case ">=":
		return recordVal >= start
	case "<":
		return recordVal < start
	case "<=":
		return recordVal <= end
	}
	return false
}

func (e *AstEvaluator) VisitGroupingExpression(expr *parser.Grouping) any {
	return expr.Expression.Accept(e).(bool)
}

func (e *AstEvaluator) VisitContainsExpression(expr *parser.Contains) any {
	content, ok := e.record["content"].(string)
	if !ok {
		return false
	}

	words := util.Tokenize(content)
	for _, w := range expr.Words {
		for _, l := range words {
			if l == strings.ToLower(w) {
				return true
			}
		}
	}

	return false
}

func (e *AstEvaluator) VisitTypeFilterExpression(expr *parser.TypeFilter) any {
	tVal, ok := e.record["type"].(string)
	if !ok {
		return false
	}

	for _, w := range expr.Words {
		if t, ok := entry.TypeFor(w); ok && string(t) == tVal {
			return true
		}
	}

	return false
}

func (e *AstEvaluator) VisitLiteralExpression(expr *parser.Literal) any {
	t := expr.Val

	switch t.TokenType {
	case lexer.NUMBER, lexer.STRING:
		return t.Literal
	}
	return nil
}

func (e *AstEvaluator) VisitTagsExpression(expr *parser.Tags) any {
	eTags := e.record["tags"].([]string)
	qTags := expr.TagList

	for _, qT := range qTags {
		if slices.Contains(eTags, qT) {
			return true
		}
	}

	return false
}

func (e *AstEvaluator) VisitIDsExpression(expr *parser.IDs) any {
	id, ok := e.record["id"].(int)
	if !ok {
		return false
	}

	return slices.Contains(expr.List, id)
}

func (e *AstEvaluator) VisitAllExpression(expr *parser.All) any {
	return true
}

type PayloadEvaluator struct{}

func EvaluatePayload(p parser.Payload) (any, any, any) {
	e := &PayloadEvaluator{}
	return p.Accept(e)
}

func (e *PayloadEvaluator) VisitRemPayload(p *parser.RemPayload) (any, any, any) {
	return p.Type, p.Tags, p.Content
}

func (e *PayloadEvaluator) VisitRevisePayload(p *parser.RevisePayload) (any, any, any) {
	return p.Type, p.Tags, nil
}
