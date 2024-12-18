package writer

import (
	"database/sql"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var LongWriterNew func(
	io.Writer,
) func(int64) IO[Void] = Curry[func(int64) int64](IntegerWriterNew)(
	func(i int64) int64 { return i },
)

var NullLongWriterNew func(
	io.Writer,
) func(
	sql.Null[int64],
) IO[Void] = Curry[func(int64) int64](NullIntegerWriterNew)(
	func(i int64) int64 { return i },
)
