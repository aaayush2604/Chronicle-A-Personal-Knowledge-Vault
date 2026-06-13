package engine

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	queryExec "chronicle/internal/query/execution"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/semantic"
)

func entryToRecord(e entry.KnowledgeEntry) queryExec.Record {
	return queryExec.Record{
		"content": e.Content,
		"date":    e.Timestamp,
		"len":     len(e.Content),
		"type":    string(e.Type),
	}
}

func (e *Engine) Query(input string) ([]entry.KnowledgeEntry, error) {
	var res []entry.KnowledgeEntry
	scanner := lexer.NewScanner(input)

	tokens := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in Query:")
	}

	err = semantic.AnalyzeSemantics(q)
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in Query:")
	}

	rootOperator := queryExec.GetExecutionRoot(q)

	eContext := &queryExec.ExecContext{
		Store:   e.store,
		Ast:     q.Expr,
		Payload: q.Payload,
	}

	if err := rootOperator.Setup(eContext); err != nil {
		return nil, errorC.Wrap(err, errorC.Execution, "Setup failed:")
	}
	defer rootOperator.Free(eContext)

	switch rootOperator.GetType() {
	case queryExec.RecallType:
		for {
			e, exhausted, err := rootOperator.Next(eContext)
			if err != nil {
				return nil, errorC.Wrap(err, errorC.Execution, "Error in Query:")
			}
			if exhausted {
				return res, nil
			}
			res = append(res, e)
		}
	case queryExec.NoteType:
		e, err := rootOperator.Write(eContext)
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Execution, "Error in Query:")
		}

		res = append(res, e)
	}

	return res, nil
}
