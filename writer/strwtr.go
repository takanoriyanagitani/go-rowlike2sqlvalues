package writer

import (
	"bytes"
	"context"
	"database/sql"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func StringWriterNew(w io.Writer) func(string) IO[Void] {
	var buf bytes.Buffer
	return func(s string) IO[Void] {
		return func(_ context.Context) (Void, error) {
			buf.Reset()
			_, _ = buf.WriteString(s) // error is always nil or panic
			_, e := w.Write(buf.Bytes())
			return Empty, e
		}
	}
}

func NullStringWriterNew(w io.Writer) func(sql.Null[string]) IO[Void] {
	var buf bytes.Buffer
	return func(s sql.Null[string]) IO[Void] {
		return func(_ context.Context) (Void, error) {
			buf.Reset()

			switch s.Valid {
			case false:
				_, _ = buf.WriteString("None") // error is always nil or panic
			default:
				// those won't return non-nil error or panic
				_, _ = buf.WriteString("Some(")
				_, _ = buf.WriteString(s.V)
				_, _ = buf.WriteString(")")
			}

			_, e := w.Write(buf.Bytes())
			return Empty, e
		}
	}
}
