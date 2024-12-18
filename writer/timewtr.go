package writer

import (
	"context"
	"database/sql"
	"io"
	"time"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func TimeWriterNew(w io.Writer) func(time.Time) IO[Void] {
	return func(s time.Time) IO[Void] {
		return func(_ context.Context) (Void, error) {
			var ts string = s.String()
			_, e := io.WriteString(w, ts)
			return Empty, e
		}
	}
}

func NullTimeWriterNew(w io.Writer) func(sql.Null[time.Time]) IO[Void] {
	return func(s sql.Null[time.Time]) IO[Void] {
		return func(_ context.Context) (Void, error) {
			if !s.Valid {
				_, e := w.Write([]byte("None"))
				return Empty, e
			}

			_, e := w.Write([]byte("Some("))
			if nil != e {
				return Empty, e
			}

			var ts string = s.V.String()
			_, e = io.WriteString(w, ts)
			if nil != e {
				return Empty, e
			}

			_, e = w.Write([]byte(")"))
			return Empty, e
		}
	}
}
