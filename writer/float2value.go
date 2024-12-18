package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func FloatToValue(p float32) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(float32) IO[Void] = pw.FloatWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableFloatToValue(p sql.Null[float32]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[float32]) IO[Void] = nw.
				FloatWriter
			return wtr(p)(ctx)
		}
	}
}
