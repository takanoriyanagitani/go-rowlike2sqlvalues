package writer

import (
	"context"
	"database/sql"
	"io"
	"strconv"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

type FloatWriterConfig[T any] struct {
	Format    byte
	Precision int
	BitSize   int
	ToDouble  func(T) float64
}

func FloatNumberWriterNew[T any](
	cfg FloatWriterConfig[T],
	w io.Writer,
) func(T) IO[Void] {
	var buf []byte
	return func(s T) IO[Void] {
		return func(_ context.Context) (Void, error) {
			buf = buf[:0]
			buf = strconv.AppendFloat(
				buf,
				cfg.ToDouble(s),
				cfg.Format,
				cfg.Precision,
				cfg.BitSize,
			)
			_, e := w.Write(buf)
			return Empty, e
		}
	}
}

func NullFloatNumberWriterNew[T any](
	cfg FloatWriterConfig[T],
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

			buf = strconv.AppendFloat(
				buf,
				cfg.ToDouble(s.V),
				cfg.Format,
				cfg.Precision,
				cfg.BitSize,
			)
			_, e = w.Write(buf)
			if nil != e {
				return Empty, e
			}

			_, e = w.Write([]byte(")"))
			return Empty, e
		}
	}
}

var FloatWriterNewDefault func(
	io.Writer,
) func(float32) IO[Void] = Curry[FloatWriterConfig[float32]](
	FloatNumberWriterNew,
)(
	FloatWriterConfig[float32]{
		Format:    'g',
		Precision: -1,
		BitSize:   32,
		ToDouble:  func(f float32) float64 { return float64(f) },
	},
)

var NullFloatWriterNewDefault func(
	io.Writer,
) func(sql.Null[float32]) IO[Void] = Curry[FloatWriterConfig[float32]](
	NullFloatNumberWriterNew,
)(
	FloatWriterConfig[float32]{
		Format:    'g',
		Precision: -1,
		BitSize:   32,
		ToDouble:  func(f float32) float64 { return float64(f) },
	},
)
