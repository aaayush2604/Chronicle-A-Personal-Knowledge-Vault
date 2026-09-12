package engine

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	"chronicle/internal/query/execution"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/semantic"
)

func entryToRecord(e entry.KnowledgeEntry) execution.Record {
	return execution.Record{
		"content": e.Content,
		"date":    e.Timestamp,
		"len":     len(e.Content),
		"type":    string(e.Type),
	}
}

func (e *Engine) Query(input string) ([]entry.KnowledgeEntry, error) {
	var res []entry.KnowledgeEntry
	scanner := lexer.NewScanner(input)

	tokens, err := scanner.ScanTokens()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in Query:")
	}

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in Query:")
	}

	err = semantic.AnalyzeSemantics(q)
	if err != nil {
		return nil, errorC.Wrap(err, errorC.Syntax, "Error in Query:")
	}

	rootOperator := execution.GetExecutionRoot(q)

	eContext := &execution.ExecContext{
		Store:   e.store,
		Ast:     q.Expr,
		Payload: q.Payload,
	}

	if err := rootOperator.Setup(eContext); err != nil {
		return nil, errorC.Wrap(err, errorC.Execution, "Setup failed:")
	}
	defer rootOperator.Free(eContext)

	switch rootOperator.GetType() {
	case execution.RecallType, execution.ForgetType:
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
	case execution.RemType:
		e, err := rootOperator.Write(eContext)
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Execution, "Error in Query:")
		}
		res = append(res, e)
	case execution.ReviseType:
		entries, err := rootOperator.Update(eContext)
		if err != nil {
			return nil, errorC.Wrap(err, errorC.Execution, "Error in Query:")
		}
		for _, e := range entries {
			res = append(res, e)
		}
	}

	return res, nil
}
