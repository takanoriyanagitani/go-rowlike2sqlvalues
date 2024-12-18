package writer

import (
	"context"
	"database/sql"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func BytesWriterNew(w io.Writer) func([]byte) IO[Void] {
	return func(s []byte) IO[Void] {
		return func(_ context.Context) (Void, error) {
			_, e := w.Write(s)
			return Empty, e
		}
	}
}

func NullBytesWriterNew(w io.Writer) func(sql.Null[[]byte]) IO[Void] {
	return func(s sql.Null[[]byte]) IO[Void] {
		return func(_ context.Context) (Void, error) {
			if !s.Valid {
				_, e := w.Write([]byte("None"))
				return Empty, e
			}

			_, e := w.Write([]byte("Some("))
			if nil != e {
				return Empty, e
			}

			_, e = w.Write(s.V)
			if nil != e {
				return Empty, e
			}

			_, e = w.Write([]byte(")"))
			return Empty, e
		}
	}
}
