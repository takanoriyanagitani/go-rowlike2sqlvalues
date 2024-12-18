package writer

import (
	"context"
	"database/sql"
	"io"

	"github.com/google/uuid"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func UuidWriterNew(w io.Writer) func(uuid.UUID) IO[Void] {
	return func(s uuid.UUID) IO[Void] {
		return func(_ context.Context) (Void, error) {
			var ts string = s.String()
			_, e := io.WriteString(w, ts)
			return Empty, e
		}
	}
}

func NullUuidWriterNew(w io.Writer) func(sql.Null[uuid.UUID]) IO[Void] {
	return func(s sql.Null[uuid.UUID]) IO[Void] {
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
