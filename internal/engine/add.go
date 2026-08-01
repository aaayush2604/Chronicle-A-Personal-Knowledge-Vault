package engine

import (
	"chronicle/internal/entry"
	"chronicle/internal/errorC"
	"chronicle/internal/query/lexer"
	"fmt"
)

func (e *Engine) AddNote(content string, t entry.EntryType) (entry.KnowledgeEntry, error) {
	ke, err := e.store.Add(content, []*lexer.Token{}, t)

	if err != nil {
		return entry.KnowledgeEntry{}, err
	}
	e.index.Build(e.store.List())

	return ke, nil
}

func (e *Engine) Delete(id int) error {
	return e.store.Delete(id)
}

func (e *Engine) ProcessDeletion(ids []int, list []entry.KnowledgeEntry) error {
	idSet := map[int]struct{}{}

	if len(ids) > 0 {
		for _, en := range list {
			idSet[en.ID] = struct{}{}
		}
		for _, id := range ids {
			if _, ok := idSet[id]; !ok {
				return errorC.New(errorC.NotFound, fmt.Sprintf("Entry with id [%d] is not part of the forget query's result", id))
			}
		}
		for _, id := range ids {
			err := e.store.Delete(id)
			if err != nil {
				return errorC.Wrap(err, errorC.Execution, fmt.Sprintf("Unable to Forget id [%d]", id))
			}
		}
		return nil
	}

	for _, en := range list {
		err := e.store.Delete(en.ID)
		if err != nil {
			return errorC.Wrap(err, errorC.Execution, fmt.Sprintf("Unable to Forget id [%d]", en.ID))
		}
	}

	return nil
}
