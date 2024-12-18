package writer

// This file is generated using any2val.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

import "database/sql"

import "github.com/google/uuid"

func UuidToValue(p uuid.UUID) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw PrimitiveWriter = vw.PrimitiveWriter
			var wtr func(uuid.UUID) IO[Void] = pw.UuidWriter
			return wtr(p)(ctx)
		}
	}
}

func NullableUuidToValue(p sql.Null[uuid.UUID]) Value {
	return func(vw ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw NullableWriter = vw.NullableWriter
			var wtr func(sql.Null[uuid.UUID]) IO[Void] = nw.
				UuidWriter
			return wtr(p)(ctx)
		}
	}
}
