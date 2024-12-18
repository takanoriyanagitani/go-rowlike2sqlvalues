package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func StringToValue(p string) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(string) IO[Void] = pw.StringWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableStringToValue(p sql.Null[string]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[string]) IO[Void] = nw.
				StringWriter
			return wtr(p)(ctx)
		}
	}
}
