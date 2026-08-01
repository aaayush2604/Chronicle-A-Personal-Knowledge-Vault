package execution

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
)

type CmdOperatorType string

const (
	RecallType CmdOperatorType = "recall"
	RemType    CmdOperatorType = "remember"
	ForgetType CmdOperatorType = "forget"
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
		if EvaluateRecall(this.ast, EntryToRecord(e)) {
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
type Cmd interface {
	GetType() CmdOperatorType
	Setup(context *ExecContext) error
	Free(context *ExecContext) error
	Next(context *ExecContext) (entry.KnowledgeEntry, bool, error)
	Write(context *ExecContext) (entry.KnowledgeEntry, error)
}

type Recall struct {
	cmdType CmdOperatorType
	ast     parser.Expr
	child   Operator
}

func NewRecall(tree parser.Expr) *Recall {
	return &Recall{
		cmdType: RecallType,
		ast:     tree,
	}
}

func (this *Recall) GetType() CmdOperatorType {
	return this.cmdType
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

func (this *Recall) Write(context *ExecContext) (entry.KnowledgeEntry, error) {
	return entry.KnowledgeEntry{}, nil
}

type Remember struct {
	cmdType CmdOperatorType
	payload parser.Payload
}

func NewRemember(p parser.Payload) *Remember {
	return &Remember{
		cmdType: RemType,
		payload: p,
	}
}

func (this *Remember) GetType() CmdOperatorType {
	return this.cmdType
}

func (this *Remember) Setup(context *ExecContext) error {
	return nil
}

func (this *Remember) Free(context *ExecContext) error {
	this.payload = nil
	return nil
}

func (this *Remember) Next(context *ExecContext) (entry.KnowledgeEntry, bool, error) {
	return entry.KnowledgeEntry{}, false, nil
}

func (this *Remember) Write(context *ExecContext) (entry.KnowledgeEntry, error) {
	eType, tags, content := EvaluatePayload(this.payload)

	e, err := context.Store.Add(TokensToString(content.([]*lexer.Token)), tags.([]*lexer.Token), eType.(entry.EntryType))
	if err != nil {
		return entry.KnowledgeEntry{}, errorC.Wrap(err, errorC.Execution, "Error in Adding Entry to Store")
	}

	return e, nil
}

type Forget struct {
	cmdType CmdOperatorType
	ast     parser.Expr
	child   Operator
}

func NewForget(tree parser.Expr) *Forget {
	return &Forget{
		cmdType: ForgetType,
		ast:     tree,
	}
}

func (this *Forget) GetType() CmdOperatorType {
	return this.cmdType
}

func (this *Forget) Next(context *ExecContext) (entry.KnowledgeEntry, bool, error) {
	e, exhausted, err := this.child.next(context)
	if err != nil {
		return entry.KnowledgeEntry{}, false, err
	}
	if exhausted {
		return entry.KnowledgeEntry{}, true, nil
	}
	return e, false, nil
}

func (this *Forget) Setup(context *ExecContext) error {
	logScan := NewLogScan()
	filter := NewFilter(logScan, this.ast)
	this.child = filter
	err := this.child.setup(context)
	if err != nil {
		return err
	}
	return nil
}

func (this *Forget) Free(context *ExecContext) error {
	err := this.child.free(context)
	if err != nil {
		return err
	}
	return nil
}

func (this *Forget) Write(context *ExecContext) (entry.KnowledgeEntry, error) {
	return entry.KnowledgeEntry{}, nil
}

func GetExecutionRoot(expr *parser.Query) Cmd {
	switch expr.Command {
	case parser.RecallCommand:
		return NewRecall(expr.Expr)
	case parser.RemCommand:
		return NewRemember(expr.Payload)
	case parser.ForgetCommand:
		return NewForget(expr.Expr)
	default:
		return nil
	}
}
