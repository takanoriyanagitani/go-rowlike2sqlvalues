package pgcopy2vals

// This file is generated using toval.tmpl. NEVER EDIT.

import (
	"context"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

func (p PgColumn) ToValueDouble() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			nb, e := p.ToNullableDouble()
			if nil != e {
				return Empty, e
			}
			if !nb.Valid {
				return Empty, ErrUnexpectedNull
			}
			return pw.DoubleWriter(nb.V)(ctx)
		}
	}
}

func (p PgColumn) ToValueNullDouble() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw sw.NullableWriter = vw.NullableWriter

			nb, e := p.ToNullableDouble()
			if nil != e {
				return Empty, e
			}

			return nw.DoubleWriter(nb)(ctx)
		}
	}
}
