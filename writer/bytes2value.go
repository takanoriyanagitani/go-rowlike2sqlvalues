package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

func BytesToValue(p []byte) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func([]byte) IO[Void] = pw.BytesWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableBytesToValue(p sql.Null[[]byte]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[[]byte]) IO[Void] = nw.
				BytesWriter
			return wtr(p)(ctx)
		}
	}
}
