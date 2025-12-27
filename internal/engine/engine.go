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
