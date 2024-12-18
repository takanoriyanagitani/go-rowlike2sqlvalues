package writer

import (
	"context"
	"database/sql"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var BoolTrue []byte = []byte("true")
var BoolFalse []byte = []byte("false")

func BoolToBytes(b bool) []byte {
	switch b {
	case true:
		return BoolTrue
	default:
		return BoolFalse
	}
}

func BooleanWriterNew(w io.Writer) func(bool) IO[Void] {
	return func(s bool) IO[Void] {
		return func(_ context.Context) (Void, error) {
			var b []byte = BoolToBytes(s)
			_, e := w.Write(b)
			return Empty, e
		}
	}
}

func NullBooleanWriterNew(w io.Writer) func(sql.Null[bool]) IO[Void] {
	return func(s sql.Null[bool]) IO[Void] {
		return func(_ context.Context) (Void, error) {
			if !s.Valid {
				_, e := w.Write([]byte("None"))
				return Empty, e
			}

			_, e := w.Write([]byte("Some("))
			if nil != e {
				return Empty, e
			}

			var b []byte = BoolToBytes(s.V)
			_, e = w.Write(b)
			if nil != e {
				return Empty, e
			}

			_, e = w.Write([]byte(")"))
			return Empty, e
		}
	}
}
