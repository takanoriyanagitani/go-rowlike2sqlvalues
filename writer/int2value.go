package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func IntToValue(p int32) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(int32) IO[Void] = pw.IntWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableIntToValue(p sql.Null[int32]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[int32]) IO[Void] = nw.
				IntWriter
			return wtr(p)(ctx)
		}
	}
}
