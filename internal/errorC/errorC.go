package errorC

import (
	"errors"
	"strings"
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func New(kind Kind, msg string) *Error {
	return &Error{
		Kind:    kind,
		Message: msg,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ": " + e.Message
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Wrap(err error, kind Kind, msg string) *Error {
	return &Error{
		Kind:    kind,
		Message: msg,
		Err:     err,
	}
}

func GetKind(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return Internal
}

func FormatError(err error) string {
	if err == nil {
		return ""
	}

	var lines []string
	current := err
	indent := 0

	for current != nil {
		if e, ok := current.(*Error); ok {
			lines = append(lines, strings.Repeat("\t", indent)+e.Message)
			current = e.Unwrap()
		} else {
			lines = append(lines, strings.Repeat("\t", indent)+current.Error())
			break
		}
		indent++
	}

	return strings.Join(lines, "\n")
}
