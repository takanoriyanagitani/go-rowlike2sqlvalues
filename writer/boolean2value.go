package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func BooleanToValue(p bool) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(bool) IO[Void] = pw.BooleanWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableBooleanToValue(p sql.Null[bool]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[bool]) IO[Void] = nw.
				BooleanWriter
			return wtr(p)(ctx)
		}
	}
}
