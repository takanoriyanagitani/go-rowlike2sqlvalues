package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func DoubleToValue(p float64) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(float64) IO[Void] = pw.DoubleWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableDoubleToValue(p sql.Null[float64]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[float64]) IO[Void] = nw.
				DoubleWriter
			return wtr(p)(ctx)
		}
	}
}
