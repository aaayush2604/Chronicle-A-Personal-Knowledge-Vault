package errorC

type Kind string

const (
	Validation Kind = "validation"
	NotFound   Kind = "not_found"
	Internal   Kind = "internal"
	Syntax     Kind = "syntax"
)
