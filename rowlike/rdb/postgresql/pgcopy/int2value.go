package pgcopy2vals

// This file is generated using toval.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

func (p PgColumn) ToValueInt() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			nb, e := p.ToNullableInt()
			if nil != e {
				return Empty, e
			}
			if !nb.Valid {
				return Empty, ErrUnexpectedNull
			}
			return pw.IntWriter(nb.V)(ctx)
		}
	}
}

func (p PgColumn) ToValueNullInt() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw sw.NullableWriter = vw.NullableWriter

			nb, e := p.ToNullableInt()
			if nil != e {
				return Empty, e
			}

			return nw.IntWriter(nb)(ctx)
		}
	}
}
