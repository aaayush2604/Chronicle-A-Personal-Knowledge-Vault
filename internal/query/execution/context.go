package execution

import (
	"chronicle/internal/query/parser"
	"chronicle/internal/store"
)

type ExecContext struct {
	Store *store.Store
	Ast   parser.Expr
}
