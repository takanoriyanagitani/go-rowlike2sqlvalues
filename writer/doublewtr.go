package writer

import (
	"database/sql"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var DoubleWriterNewDefault func(
	io.Writer,
) func(float64) IO[Void] = Curry[FloatWriterConfig[float64]](
	FloatNumberWriterNew,
)(
	FloatWriterConfig[float64]{
		Format:    'g',
		Precision: -1,
		BitSize:   64,
		ToDouble:  func(f float64) float64 { return f },
	},
)

var NullDoubleWriterNewDefault func(
	io.Writer,
) func(sql.Null[float64]) IO[Void] = Curry[FloatWriterConfig[float64]](
	NullFloatNumberWriterNew,
)(
	FloatWriterConfig[float64]{
		Format:    'g',
		Precision: -1,
		BitSize:   64,
		ToDouble:  func(f float64) float64 { return f },
	},
)
