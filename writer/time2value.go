package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

import "time"

func TimeToValue(p time.Time) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(time.Time) IO[Void] = pw.TimeWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableTimeToValue(p sql.Null[time.Time]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[time.Time]) IO[Void] = nw.
				TimeWriter
			return wtr(p)(ctx)
		}
	}
}
