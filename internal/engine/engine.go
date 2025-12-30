package engine

import (
	"chronicle/internal/index"
	"chronicle/internal/store"
)

type Engine struct {
	store *store.Store
	index *index.Index
}

func New(store *store.Store, index *index.Index) *Engine {
	index.Build(store.List())

	return &Engine{
		store: store,
		index: index,
	}
}

func (e *Engine) PrintIndex() {
	e.index.PrintIndex()
}

func (e *Engine) CheckDelete(id int) (bool, store.DeletionInfo) {
	if deleted, ts := e.store.CheckDelete(id); deleted {
		return deleted, ts
	}
	return false, store.DeletionInfo{}
}
