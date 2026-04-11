package engine

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	queryExec "chronicle/internal/query/execution"
	"chronicle/internal/query/lexer"
	"chronicle/internal/query/parser"
	"chronicle/internal/query/semantic"
	"fmt"
)

func entryToRecord(e entry.KnowledgeEntry) queryExec.Record {
	return queryExec.Record{
		"content": e.Content,
		"date":    e.Timestamp,
		"len":     len(e.Content),
		"type":    string(e.Type),
	}
}

func (e *Engine) Query(input string) []entry.KnowledgeEntry {
	entries := e.store.List()
	var res []entry.KnowledgeEntry

	scanner := lexer.NewScanner(input)

	tokens := scanner.ScanTokens()

	p := parser.NewParser(tokens)

	q, err := p.Parse()
	if err != nil {
		fmt.Println(errorC.FormatError(err))
	}

	err = semantic.AnalyzeSemantics(q)
	if err != nil {
		errorString := errorC.FormatError(err)
		fmt.Println(errorString)
	}

	for _, e := range entries {
		r := entryToRecord(e)

		if queryExec.Evaluate(q.Expr, r) {
			res = append(res, e)
		}
	}

	return res
}
