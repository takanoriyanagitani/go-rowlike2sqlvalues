package pgcopy2vals

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"iter"
	"os"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

type PgRow []PgColumn

func MapGetOrInsert[K comparable, V any](
	m map[K]V,
	key K,
	onMissing func() V,
) V {
	val, found := m[key]
	if found {
		return val
	}

	var neo V = onMissing()
	m[key] = neo
	return m[key]
}

func ReaderToPgRows(
	rdr io.Reader,
) iter.Seq2[PgRow, error] {
	return func(yield func(PgRow, error) bool) {
		var br io.Reader = bufio.NewReader(rdr)

		var colcntbuf [2]byte

		var colsizbuf [4]byte

		var row PgRow

		var bufs map[int16]*bytes.Buffer = map[int16]*bytes.Buffer{}

		for {
			_, e := io.ReadFull(br, colcntbuf[:])
			if nil != e {
				yield(nil, e)
				return
			}

			var colcntu uint16 = binary.BigEndian.Uint16(colcntbuf[:])
			var colcnt int16 = int16(colcntu)
			if -1 == colcnt {
				return
			}

			var ix int16
			clear(row)
			row = row[:0]
			for ix = 0; ix < colcnt; ix++ {
				_, e := io.ReadFull(br, colsizbuf[:])
				if nil != e {
					yield(nil, e)
					return
				}

				var colsizu uint32 = binary.BigEndian.Uint32(colsizbuf[:])
				var colsiz int32 = int32(colsizu)
				var isNullCol bool = colsiz < 0
				var col PgColumn
				col.Size = colsiz
				col.Content = nil

				if !isNullCol {
					var buf *bytes.Buffer = MapGetOrInsert(
						bufs,
						ix,
						func() *bytes.Buffer {
							var b bytes.Buffer
							return &b
						},
					)
					buf.Reset()
					if 0 < col.Size {
						limited := &io.LimitedReader{
							R: br,
							N: int64(col.Size),
						}
						_, e := io.Copy(buf, limited)
						if nil != e {
							yield(nil, e)
							return
						}
					}
					col.Content = buf.Bytes()
				}

				row = append(row, col)
			}

			if !yield(row, nil) {
				return
			}
		}
	}
}

func StdinToPgRows() iter.Seq2[PgRow, error] { return ReaderToPgRows(os.Stdin) }

var PgRowsFromStdin IO[iter.Seq2[PgRow, error]] = OfFn(StdinToPgRows)
