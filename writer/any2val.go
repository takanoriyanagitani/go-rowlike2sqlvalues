package writer

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidValue error = errors.New("invalid value")
)

func ConvertToNullable[T any](t *T) sql.Null[T] {
	var ret sql.Null[T]
	switch t {
	case nil:
		return ret
	default:
		ret.V = *t
	}
	return ret
}

//go:generate go run internal/gen/a2vsw/main.go
type AnyToValue func(any) Value
