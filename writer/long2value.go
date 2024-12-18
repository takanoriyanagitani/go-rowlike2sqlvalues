package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func LongToValue(p int64) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(int64) IO[Void] = pw.LongWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableLongToValue(p sql.Null[int64]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[int64]) IO[Void] = nw.
				LongWriter
			return wtr(p)(ctx)
		}
	}
}
