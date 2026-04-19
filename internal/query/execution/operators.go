package execution

import (
	"chronicle/internal/entry"
	"chronicle/internal/query/parser"
)

// execution pipeline operators
type Operator interface {
	next(context *ExecContext) (entry.KnowledgeEntry, bool, error)
	setup(context *ExecContext) error
	free(context *ExecContext) error
}

type LogScan struct {
	currentRow int
	entries    []entry.KnowledgeEntry
}

func NewLogScan() *LogScan {
	return &LogScan{
		currentRow: 0,
	}
}

func (this *LogScan) next(context *ExecContext) (entry.KnowledgeEntry, bool, error) {
	if this.currentRow >= len(this.entries) {
		return entry.KnowledgeEntry{}, true, nil
	}

	e := this.entries[this.currentRow]
	this.currentRow++
	return e, false, nil
}

func (this *LogScan) setup(context *ExecContext) error {
	this.currentRow = 0
	this.entries = context.Store.List()
	return nil
}

func (this *LogScan) free(context *ExecContext) error {
	this.currentRow = 0
	this.entries = nil
	return nil
}

type Filter struct {
	child Operator
	ast   parser.Expr
}

func NewFilter(op Operator, tree parser.Expr) *Filter {
	return &Filter{
		child: op,
		ast:   tree,
	}
}

func (this *Filter) next(context *ExecContext) (entry.KnowledgeEntry, bool, error) {
	for {
		e, exhausted, err := this.child.next(context)
		if err != nil {
			return entry.KnowledgeEntry{}, false, err
		}
		if exhausted {
			return entry.KnowledgeEntry{}, true, nil
		}
		if Evaluate(this.ast, EntryToRecord(e)) {
			return e, false, nil
		}
	}
}

func (this *Filter) setup(context *ExecContext) error {
	err := this.child.setup(context)
	if err != nil {
		return err
	}
	return nil
}

func (this *Filter) free(context *ExecContext) error {
	err := this.child.free(context)
	if err != nil {
		return err
	}
	this.child = nil
	this.ast = nil
	return nil
}

// root operators
type Command interface {
	Setup(context *ExecContext) error
	Free(context *ExecContext) error
	Next(context *ExecContext) (entry.KnowledgeEntry, bool, error)
}

type Recall struct {
	ast   parser.Expr
	child Operator
}

func NewRecall(tree parser.Expr) *Recall {
	return &Recall{
		ast: tree,
	}
}

func (this *Recall) Next(context *ExecContext) (entry.KnowledgeEntry, bool, error) {
	e, exhausted, err := this.child.next(context)
	if err != nil {
		return entry.KnowledgeEntry{}, false, err
	}
	if exhausted {
		return entry.KnowledgeEntry{}, true, nil
	}
	return e, false, nil
}

func (this *Recall) Setup(context *ExecContext) error {
	logScan := NewLogScan()
	filter := NewFilter(logScan, this.ast)
	this.child = filter
	err := this.child.setup(context)
	if err != nil {
		return err
	}
	return nil
}

func (this *Recall) Free(context *ExecContext) error {
	err := this.child.free(context)
	if err != nil {
		return err
	}
	return nil
}

func GetExecutionRoot(expr *parser.Query) Command {
	switch expr.Command {
	case parser.RecallCommand:
		return NewRecall(expr.Expr)
	default:
		return nil
	}
}
