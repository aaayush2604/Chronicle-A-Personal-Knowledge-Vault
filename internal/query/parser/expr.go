package parser

type CommandType string

const (
	RecallCommand CommandType = "recall"
)

type ExprVisitor interface {
	VisitLogicalExpression(*Logical) any
	VisitComparisonExpression(*Comparison) any
	VisitContainsExpression(*Contains) any
	VisitTypeFilterExpression(*TypeFilter) any
	VisitGroupingExpression(*Grouping) any
}

type Expr interface {
	exprNode()
	Accept(v ExprVisitor) any
}

type Query struct {
	Command CommandType
	Expr    Expr
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l *Logical) exprNode() {}

func (l *Logical) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpression(l)
}

type Comparison struct {
	Field    string
	Operator Token
	Value    interface{}
}

func (c *Comparison) exprNode() {}

func (c *Comparison) Accept(v ExprVisitor) any {
	return v.VisitComparisonExpression(c)
}

type Contains struct {
	Words []string
}

func (c *Contains) exprNode() {}

func (c *Contains) Accept(v ExprVisitor) any {
	return v.VisitContainsExpression(c)
}

type TypeFilter struct {
	Words []string
}

func (t *TypeFilter) exprNode() {}

func (t *TypeFilter) Accept(v ExprVisitor) any {
	return v.VisitTypeFilterExpression(t)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) exprNode() {}

func (g *Grouping) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpression(g)
}
