package parser

import (
	"chronicle/internal/entry"
	lexer "chronicle/internal/query/lexer"
)

type CommandType string

const (
	RecallCommand CommandType = "recall"
	RemCommand    CommandType = "remember"
	ForgetCommand CommandType = "forget"
)

type ExprVisitor interface {
	VisitLogicalExpression(*Logical) any
	VisitComparisonExpression(*Comparison) any
	VisitContainsExpression(*Contains) any
	VisitTypeFilterExpression(*TypeFilter) any
	VisitGroupingExpression(*Grouping) any
	VisitLiteralExpression(*Literal) any
	VisitTagsExpression(*Tags) any
	VisitAllExpression(*All) any
}

type PayloadVisitor interface {
	VisitRemPayload(*RemPayload) (any, any, any)
}

type Expr interface {
	exprNode()
	Accept(v ExprVisitor) any
}

type Payload interface {
	payloadNode()
	Accept(v PayloadVisitor) (any, any, any)
}

type Query struct {
	Command CommandType
	Expr    Expr
	Payload Payload
}

type Logical struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func NewLogical(left Expr, Op *lexer.Token, right Expr) *Logical {
	return &Logical{
		Left:     left,
		Operator: Op,
		Right:    right,
	}
}

func (l *Logical) exprNode() {}

func (l *Logical) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpression(l)
}

type Comparison struct {
	Field    *lexer.Token
	Operator *lexer.Token
	Value    *Literal
}

func NewComparison(field *lexer.Token, Op *lexer.Token, v *Literal) *Comparison {
	return &Comparison{
		Field:    field,
		Operator: Op,
		Value:    v,
	}
}

func (c *Comparison) exprNode() {}

func (c *Comparison) Accept(v ExprVisitor) any {
	return v.VisitComparisonExpression(c)
}

type Contains struct {
	Words []string
}

func NewContains(words []string) *Contains {
	return &Contains{
		Words: words,
	}
}

func (c *Contains) exprNode() {}

func (c *Contains) Accept(v ExprVisitor) any {
	return v.VisitContainsExpression(c)
}

type TypeFilter struct {
	Words []string
}

func NewTypeFilter(words []string) *TypeFilter {
	return &TypeFilter{
		Words: words,
	}
}

func (t *TypeFilter) exprNode() {}

func (t *TypeFilter) Accept(v ExprVisitor) any {
	return v.VisitTypeFilterExpression(t)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

func (g *Grouping) exprNode() {}

func (g *Grouping) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpression(g)
}

type Literal struct {
	Val *lexer.Token
}

func NewLiteral(val *lexer.Token) *Literal {
	return &Literal{
		Val: val,
	}
}

func (l *Literal) exprNode() {}

func (l *Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpression(l)
}

type Tags struct {
	TagList []string
}

func NewTags(tags []string) *Tags {
	return &Tags{
		TagList: tags,
	}
}

func (t *Tags) exprNode() {}

func (t *Tags) Accept(v ExprVisitor) any {
	return v.VisitTagsExpression(t)
}

type All struct{}

func (a *All) exprNode() {}

func (a *All) Accept(v ExprVisitor) any {
	return v.VisitAllExpression(a)
}

type RemPayload struct {
	Type    entry.EntryType
	Tags    []*lexer.Token
	Content []*lexer.Token
}

func (n *RemPayload) payloadNode() {}

func (n *RemPayload) Accept(v PayloadVisitor) (any, any, any) {
	return v.VisitRemPayload(n)
}
