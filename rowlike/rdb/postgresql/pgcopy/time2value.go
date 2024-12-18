package pgcopy2vals

// This file is generated using toval.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

func (p PgColumn) ToValueTime() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			nb, e := p.ToNullableTime()
			if nil != e {
				return Empty, e
			}
			if !nb.Valid {
				return Empty, ErrUnexpectedNull
			}
			return pw.TimeWriter(nb.V)(ctx)
		}
	}
}

func (p PgColumn) ToValueNullTime() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw sw.NullableWriter = vw.NullableWriter

			nb, e := p.ToNullableTime()
			if nil != e {
				return Empty, e
			}

			return nw.TimeWriter(nb)(ctx)
		}
	}
}
