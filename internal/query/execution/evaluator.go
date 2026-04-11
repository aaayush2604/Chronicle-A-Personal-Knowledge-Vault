package execution

import (
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/util"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Record map[string]any

type Evaluator struct {
	record Record
}

func Evaluate(expr parser.Expr, record Record) bool {
	e := &Evaluator{record: record}
	return expr.Accept(e).(bool)
}

func (e *Evaluator) VisitLogicalExpression(expr *parser.Logical) any {
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

func (e *Evaluator) VisitComparisonExpression(expr *parser.Comparison) any {
	inputField := expr.Field.Lexeme
	field := inputField

	if inputField == "time" {
		field = "date"
	}

	recordVal := e.record[field]
	queryVal := expr.Value.Accept(e)

	switch inputField {
	case "date":
		adjustDate(&recordVal, &queryVal)
	case "time":
		adjustTime(&recordVal, &queryVal)
	}

	switch expr.Operator.Lexeme {
	case "<":
		return recordVal.(int) < queryVal.(int)
	case ">":
		return recordVal.(int) > queryVal.(int)
	case "=":
		return recordVal.(int) == queryVal.(int)
	case ">=":
		return recordVal.(int) >= queryVal.(int)
	case "<=":
		return recordVal.(int) <= queryVal.(int)
	case "!=":
		return recordVal.(int) != queryVal.(int)
	}

	return false
}

func (e *Evaluator) VisitGroupingExpression(expr *parser.Grouping) any {
	return expr.Accept(e).(bool)
}

func (e *Evaluator) VisitContainsExpression(expr *parser.Contains) any {
	content := e.record["content"].(string)
	flag := false

	words := util.Tokenize(content)
	for _, w := range expr.Words {
		for _, l := range words {
			if l == strings.ToLower(w) {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
		flag = false
	}

	return true
}

func (e *Evaluator) VisitTypeFilterExpression(expr *parser.TypeFilter) any {
	tVal := e.record["type"].(string)

	for _, t := range expr.Words {
		switch t {
		case "n":
			t = "note"
		case "q":
			t = "question"
		case "l":
			t = "learning"
		case "i":
			t = "idea"
		}
	}

	return slices.Contains(expr.Words, tVal)
}

func (e *Evaluator) VisitLiteralExpression(expr *parser.Literal) any {
	t := expr.Val

	switch t.TokenType {
	case lexer.NUMBER:
		val, _ := strconv.Atoi(t.Lexeme)
		return val
	case lexer.STRING:
		return t.Literal
	}
	return nil
}

func adjustDate(rVal *any, qVal *any) {
	t := (*rVal).(time.Time)
	q := (*qVal).(string)

	q = strings.ReplaceAll(q, "/", "-")

	parts := strings.Split(q, "-")

	var rInt, qInt int

	switch len(parts) {
	case 1:
		year, _ := strconv.Atoi(parts[0])
		qInt = year
		rInt = t.Year()

	case 2:
		month, _ := strconv.Atoi(parts[0])
		year, _ := strconv.Atoi(parts[1])

		qInt = year*100 + month
		rInt = t.Year()*100 + int(t.Month())

	case 3:
		day, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		year, _ := strconv.Atoi(parts[2])

		qInt = year*10000 + month*100 + day
		rInt = t.Year()*10000 + int(t.Month())*100 + t.Day()
	}

	*rVal = rInt
	*qVal = qInt
}

func adjustTime(rVal *any, qVal *any) {
	t := (*rVal).(time.Time)
	q := (*qVal).(string)

	q = strings.ToLower(strings.TrimSpace(q))

	isPM := strings.Contains(q, "pm")
	isAM := strings.Contains(q, "am")

	q = strings.ReplaceAll(q, "am", "")
	q = strings.ReplaceAll(q, "pm", "")
	q = strings.TrimSpace(q)

	parts := strings.Split(q, ":")

	var hour, min, sec int

	if len(parts) >= 1 {
		hour, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		min, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		sec, _ = strconv.Atoi(parts[2])
	}

	if isPM && hour != 12 {
		hour += 12
	}
	if isAM && hour == 12 {
		hour = 0
	}

	var rInt, qInt int

	switch len(parts) {
	case 1:
		qInt = hour
		rInt = t.Hour()

	case 2:
		qInt = hour*100 + min
		rInt = t.Hour()*100 + t.Minute()

	case 3:
		qInt = hour*10000 + min*100 + sec
		rInt = t.Hour()*10000 + t.Minute()*100 + t.Second()
	}

	*rVal = rInt
	*qVal = qInt
}
