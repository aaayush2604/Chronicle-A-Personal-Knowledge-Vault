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
	VisitLiteralExpression(*Literal) any
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

func NewLogical(left Expr, Op Token, right Expr) *Logical {
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
	Field    string
	Operator Token
	Value    Literal
}

func NewComparison(field string, Op Token, v Literal) *Comparison {
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
	Val interface{}
}

func NewLiteral(val any) *Literal {
	return &Literal{
		Val: val,
	}
}

func (l *Literal) exprNode() {}

func (l *Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpression(l)
}
