package pgcopy2vals

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"iter"
	"maps"
	"os"

	rs "github.com/takanoriyanagitani/go-rowlike2sqlvalues"
	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

type ColumnInfo struct {
	Name string
	rs.PrimitiveType
}

func ColumInfoNew(name string, typeString string) (ColumnInfo, error) {
	typ, e := rs.StringToType(typeString)
	return ColumnInfo{
		Name:          name,
		PrimitiveType: typ,
	}, e
}

type ColumnInfoRaw struct {
	Name                string `json:"name"`
	PrimitiveTypeString string `json:"type"`
}

func (r ColumnInfoRaw) ToColumnInfo() (ColumnInfo, error) {
	return ColumInfoNew(r.Name, r.PrimitiveTypeString)
}

func JsonsToColumnInfo(
	jsons io.Reader,
) iter.Seq2[ColumnInfo, error] {
	return func(yield func(ColumnInfo, error) bool) {
		var s *bufio.Scanner = bufio.NewScanner(jsons)

		for s.Scan() {
			var line []byte = s.Bytes()
			var buf ColumnInfoRaw
			e := json.Unmarshal(line, &buf)
			ci, ce := buf.ToColumnInfo()
			if !yield(ci, errors.Join(e, ce)) {
				return
			}
		}
	}
}

func JsonsFilenameToColumnInfo(
	filename string,
) iter.Seq2[ColumnInfo, error] {
	return func(yield func(ColumnInfo, error) bool) {
		f, e := os.Open(filename)
		if nil != e {
			yield(ColumnInfo{}, e)
			return
		}
		defer f.Close()

		var pairs iter.Seq2[ColumnInfo, error] = JsonsToColumnInfo(f)
		for col, e := range pairs {
			if !yield(col, e) {
				return
			}
		}
	}
}

func ColumnInfoMapFromIter(
	i iter.Seq2[ColumnInfo, error],
) map[int16]ColumnInfo {
	var pairs iter.Seq2[int16, ColumnInfo] = func(
		yield func(int16, ColumnInfo) bool,
	) {
		var ix int16
		for ci, e := range i {
			if nil == e {
				yield(ix, ci)
			}
			ix += 1
		}
	}
	return maps.Collect(pairs)
}

var JsonsFilenameToColumnInfoMap func(string) map[int16]ColumnInfo = Compose(
	JsonsFilenameToColumnInfo,
	ColumnInfoMapFromIter,
)

func FilenameToColumnInfoMap(filename string) IO[map[int16]ColumnInfo] {
	return OfFn(func() map[int16]ColumnInfo {
		return JsonsFilenameToColumnInfoMap(filename)
	})
}
