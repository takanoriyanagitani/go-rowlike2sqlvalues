package writer

import (
	"context"
	"database/sql"
	"io"
	"strconv"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func IntegerWriterNew[T any](
	integer2long func(T) int64,
	w io.Writer,
) func(T) IO[Void] {
	var buf []byte
	return func(s T) IO[Void] {
		return func(_ context.Context) (Void, error) {
			buf = buf[:0]
			buf = strconv.AppendInt(buf, integer2long(s), 10)
			_, e := w.Write(buf)
			return Empty, e
		}
	}
}

func NullIntegerWriterNew[T any](
	integer2long func(T) int64,
	w io.Writer,
) func(sql.Null[T]) IO[Void] {
	var buf []byte
	return func(s sql.Null[T]) IO[Void] {
		return func(_ context.Context) (Void, error) {
			buf = buf[:0]

			if !s.Valid {
				_, e := w.Write([]byte("None"))
				return Empty, e
			}

			_, e := w.Write([]byte("Some("))
			if nil != e {
				return Empty, e
			}

			buf = strconv.AppendInt(buf, integer2long(s.V), 10)
			_, e = w.Write(buf)
			if nil != e {
				return Empty, e
			}

			_, e = w.Write([]byte(")"))
			return Empty, e
		}
	}
}

var IntWriterNew func(
	io.Writer,
) func(int32) IO[Void] = Curry[func(int32) int64](IntegerWriterNew)(
	func(i int32) int64 { return int64(i) },
)

var NullIntWriterNew func(
	io.Writer,
) func(
	sql.Null[int32],
) IO[Void] = Curry[func(int32) int64](NullIntegerWriterNew)(
	func(i int32) int64 { return int64(i) },
)
