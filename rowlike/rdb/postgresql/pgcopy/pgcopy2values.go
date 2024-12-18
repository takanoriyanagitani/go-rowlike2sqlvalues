package pgcopy2vals

import (
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"maps"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	rs "github.com/takanoriyanagitani/go-rowlike2sqlvalues"
	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

var (
	ErrColumnNameNotFound error = errors.New("column name not found")

	ErrInvalidDword error = errors.New("invalid dword")
	ErrInvalidQword error = errors.New("invalid qword")

	ErrInvalidBool error = errors.New("invalid bool")

	ErrUnexpectedNull error = errors.New("unexpected null")

	ErrConvGenMissing error = errors.New("converter creator missing")

	ErrConverterMissing error = errors.New("converter missing")

	ErrInvalidUuid error = errors.New("invalid uuid")
)

//go:generate go run ./internal/gen/toval/main.go Int
//go:generate go run ./internal/gen/toval/main.go Long
//go:generate go run ./internal/gen/toval/main.go Float
//go:generate go run ./internal/gen/toval/main.go Double
//go:generate go run ./internal/gen/toval/main.go Boolean
//go:generate go run ./internal/gen/toval/main.go Time
//go:generate go run ./internal/gen/toval/main.go Uuid
//go:generate gofmt -s -w .
type PgColumn struct {
	Size    int32
	Content []byte
}

func (p PgColumn) String() string {
	var sz int32 = p.Size
	var content []byte = p.Content
	return fmt.Sprintf(
		"PgColumn(Size=%v, Content=%v)", sz, content,
	)
}

type PgColumnToValue func(PgColumn) IO[sw.Value]

type StringChecker func(string) error

var StringCheckerTrusted StringChecker = func(_ string) error { return nil }

type PgColumnConverterConfig struct{ StringChecker }

var PgColumnConverterConfigDefault PgColumnConverterConfig = PgColumnConverterConfig{
	StringChecker: StringCheckerTrusted,
}

type ConfigToConverter func(PgColumnConverterConfig) PgColumnToValue

//go:generate go run ./internal/gen/type2convmap/main.go
//go:generate gofmt -s -w .
type TypeToConverterMap map[rs.PrimitiveType]ConfigToConverter

type TypeToConfigToConverter func(rs.PrimitiveType) (ConfigToConverter, error)

var TypeToConvGen TypeToConfigToConverter = rs.
	GetValueByKeyFromMap[rs.PrimitiveType, ConfigToConverter](
	func(typ rs.PrimitiveType) error {
		return fmt.Errorf("%w: %v", ErrConvGenMissing, typ)
	},
)(type2convGenMap)

type IndexToType func(int16) rs.PrimitiveType

func ColumnIndexToColumnToValueNew(
	typ2cfg2conv TypeToConfigToConverter,
	cfg PgColumnConverterConfig,
	indices []int16,
	index2type IndexToType,
) func(int16, PgColumn) IO[sw.Value] {
	var ix2cfg2convPairs iter.Seq2[int16, ConfigToConverter] = func(
		yield func(int16, ConfigToConverter) bool,
	) {
		for _, ix := range indices {
			var typ rs.PrimitiveType = index2type(ix)
			cfg2cnv, e := typ2cfg2conv(typ)
			if nil == e {
				yield(ix, cfg2cnv)
			}
		}
	}

	var ix2convPairs iter.Seq2[int16, PgColumnToValue] = func(
		yield func(int16, PgColumnToValue) bool,
	) {
		for ix, cfg2conv := range ix2cfg2convPairs {
			var cnv PgColumnToValue = cfg2conv(cfg)
			yield(ix, cnv)
		}
	}
	var ix2convMap map[int16]PgColumnToValue = maps.Collect(ix2convPairs)
	var ix2conv func(
		int16,
	) (PgColumnToValue, error) = rs.
		GetValueByKeyFromMap[int16, PgColumnToValue](
		func(key int16) error {
			return fmt.Errorf("%w: %v", ErrConverterMissing, key)
		},
	)(ix2convMap)

	return func(ix int16, col PgColumn) IO[sw.Value] {
		conv, e := ix2conv(ix)
		if nil != e {
			return Err[sw.Value](e)
		}
		return conv(col)
	}
}

func MapToIndexToType(m map[int16]rs.PrimitiveType) IndexToType {
	return func(ix int16) rs.PrimitiveType {
		typ, found := m[ix]
		switch found {
		case true:
			return typ
		default:
			return rs.PrimitiveUnknown
		}
	}
}

func ColumnIndexToColumnToValueNewDefault(
	indices []int16,
	index2type IndexToType,
) func(int16, PgColumn) IO[sw.Value] {
	return ColumnIndexToColumnToValueNew(
		TypeToConvGen,
		PgColumnConverterConfigDefault,
		indices,
		index2type,
	)
}

func ColumnIndexToColumnToValueNewDefaultFromMap(
	ix2typMap map[int16]rs.PrimitiveType,
) func(int16, PgColumn) IO[sw.Value] {
	var ix2typ IndexToType = MapToIndexToType(ix2typMap)
	var i iter.Seq[int16] = maps.Keys(ix2typMap)
	var indices []int16 = slices.Collect(i)
	return ColumnIndexToColumnToValueNewDefault(
		indices,
		ix2typ,
	)
}

type ColumnIndexToCol2Val func(int16) PgColumnToValue

func (p PgColumn) IsNull() bool { return p.Size < 0 }

func (p PgColumn) ToValueString(
	checker func(string) error,
	buf *strings.Builder,
) sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			ns, e := p.ToNullableString(
				checker,
				buf,
			)
			if nil != e {
				return Empty, e
			}
			if !ns.Valid {
				return Empty, ErrUnexpectedNull
			}
			var s string = ns.V
			e = checker(s)
			if nil != e {
				return Empty, e
			}
			return pw.StringWriter(s)(ctx)
		}
	}
}

func (p PgColumn) ToValueNullString(
	checker func(string) error,
	buf *strings.Builder,
) sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw sw.NullableWriter = vw.NullableWriter
			ns, e := p.ToNullableString(
				checker,
				buf,
			)
			if nil != e {
				return Empty, e
			}
			var s string = ns.V
			e = checker(s)
			if nil != e {
				return Empty, e
			}
			return nw.StringWriter(ns)(ctx)
		}
	}
}

func (p PgColumn) ToValueNull() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			return pw.NullWriter(struct{}{})(ctx)
		}
	}
}

func (p PgColumn) ToValueBytes() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var pw sw.PrimitiveWriter = vw.PrimitiveWriter
			var raw []byte = p.ToBytes()
			if nil == raw {
				return Empty, ErrUnexpectedNull
			}
			return pw.BytesWriter(raw)(ctx)
		}
	}
}

func (p PgColumn) ToValueNullBytes() sw.Value {
	return func(vw sw.ValueWriter) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var nw sw.NullableWriter = vw.NullableWriter
			var raw []byte = p.ToBytes()
			var nb sql.Null[[]byte]
			nb.V = raw
			nb.Valid = nil != raw
			return nw.BytesWriter(nb)(ctx)
		}
	}
}

func (p PgColumn) ToNullableString(
	checker func(string) error,
	buf *strings.Builder,
) (sql.Null[string], error) {
	var ret sql.Null[string]
	buf.Reset()
	if p.IsNull() {
		return ret, nil
	}

	var sz int = int(p.Size)
	var raw []byte = p.Content[:sz]

	_, _ = buf.Write(raw) // error is always nil or OOM

	var s string = buf.String()

	ret.V = s
	e := checker(s)
	ret.Valid = nil == e

	return ret, e
}

func (p PgColumn) ToBytes() []byte {
	if p.IsNull() {
		return nil
	}
	var sz int = int(p.Size)
	return p.Content[:sz]
}

func ByteToBool(b byte) bool {
	switch b {
	case 0:
		return false
	default:
		return true
	}
}

const PgtimeOffset int64 = 10957 * 86400 * 1000 * 1000

func (p PgColumn) ToNullableTime() (sql.Null[time.Time], error) {
	lng, e := p.ToNullableLong()
	if nil != e {
		return sql.Null[time.Time]{}, e
	}
	return NullableMap(
		lng,
		func(i int64) time.Time {
			var unixtimeUs int64 = i + PgtimeOffset
			return time.UnixMicro(unixtimeUs)
		},
	), nil
}

func (p PgColumn) ToNullableUuid() (sql.Null[uuid.UUID], error) {
	var ret sql.Null[uuid.UUID]

	var b []byte = p.ToBytes()
	if 0 == len(b){
		return ret, nil
	}

	if 16 != len(b) {
		return ret, ErrInvalidUuid
	}

	copy(ret.V[:], b)
	ret.Valid = true
	return ret, nil
}

func (p PgColumn) ToNullableBoolean() (sql.Null[bool], error) {
	var ret sql.Null[bool]
	if p.IsNull() {
		return ret, nil
	}

	if 1 != len(p.Content) {
		return ret, fmt.Errorf("%w: %v", ErrInvalidBool, len(p.Content))
	}

	var b bool = ByteToBool(p.Content[0])

	ret.V = b
	ret.Valid = true
	return ret, nil
}

func (p PgColumn) ToNullableDword() (sql.Null[uint32], error) {
	var ret sql.Null[uint32]
	if p.IsNull() {
		return ret, nil
	}

	if 4 != len(p.Content) {
		return ret, fmt.Errorf("%w: %v", ErrInvalidDword, len(p.Content))
	}

	var buf [4]byte
	copy(buf[:], p.Content)

	var u uint32 = binary.BigEndian.Uint32(buf[:])
	ret.V = u
	ret.Valid = true
	return ret, nil
}

func (p PgColumn) ToNullableQword() (sql.Null[uint64], error) {
	var ret sql.Null[uint64]
	if p.IsNull() {
		return ret, nil
	}

	if 8 != len(p.Content) {
		return ret, fmt.Errorf("%w: %v", ErrInvalidQword, len(p.Content))
	}

	var buf [8]byte
	copy(buf[:], p.Content)

	var u uint64 = binary.BigEndian.Uint64(buf[:])
	ret.V = u
	ret.Valid = true
	return ret, nil
}

func NullableMap[T, U any](
	input sql.Null[T],
	convert func(T) U,
) sql.Null[U] {
	var ret sql.Null[U]
	if input.Valid {
		var u U = convert(input.V)
		ret.V = u
		ret.Valid = true
	}
	return ret
}

func (p PgColumn) ToNullableInt() (sql.Null[int32], error) {
	dword, e := p.ToNullableDword()
	if nil != e {
		return sql.Null[int32]{}, e
	}
	return NullableMap(
		dword,
		func(d uint32) int32 { return int32(d) },
	), nil
}

func (p PgColumn) ToNullableLong() (sql.Null[int64], error) {
	qword, e := p.ToNullableQword()
	if nil != e {
		return sql.Null[int64]{}, e
	}
	return NullableMap(
		qword,
		func(d uint64) int64 { return int64(d) },
	), nil
}

func (p PgColumn) ToNullableFloat() (sql.Null[float32], error) {
	dword, e := p.ToNullableDword()
	if nil != e {
		return sql.Null[float32]{}, e
	}
	return NullableMap(
		dword,
		math.Float32frombits,
	), nil
}

func (p PgColumn) ToNullableDouble() (sql.Null[float64], error) {
	qword, e := p.ToNullableQword()
	if nil != e {
		return sql.Null[float64]{}, e
	}
	return NullableMap(
		qword,
		math.Float64frombits,
	), nil
}

type ColumnIndexToColumnName func(int16) (string, error)

func PgRowsToValues(
	ctx context.Context,
	rows iter.Seq2[PgRow, error],
	ix2name ColumnIndexToColumnName,
	col2val func(int16, PgColumn) IO[sw.Value],
) iter.Seq2[map[string]sw.Value, error] {
	return func(yield func(map[string]sw.Value, error) bool) {
		var buf map[string]sw.Value = map[string]sw.Value{}
		for row, e := range rows {
			if nil == e {
				for i, col := range row {
					var ix int16 = int16(i)

					name, e := ix2name(ix)
					if nil != e {
						yield(buf, e)
						return
					}

					v, e := col2val(ix, col)(ctx)
					if nil != e {
						yield(buf, e)
						return
					}
					buf[name] = v
				}
			}

			if !yield(buf, e) {
				return
			}
		}
	}
}

func PgRowsToValuesFromColumnMap(
	ctx context.Context,
	rows iter.Seq2[PgRow, error],
	colmap map[int16]ColumnInfo,
) iter.Seq2[map[string]sw.Value, error] {
	var ix2names iter.Seq2[int16, string] = func(
		yield func(int16, string) bool,
	) {
		var ix2info iter.Seq2[int16, ColumnInfo] = maps.All(colmap)
		for ix, info := range ix2info {
			var name string = info.Name
			yield(ix, name)
		}
	}
	var ix2nameMap map[int16]string = maps.Collect(ix2names)
	var ix2name ColumnIndexToColumnName = rs.
		GetValueByKeyFromMap[int16, string](
		func(ix int16) error {
			return fmt.Errorf("%w: %v", ErrColumnNameNotFound, ix)
		},
	)(ix2nameMap)

	var ix2typs iter.Seq2[int16, rs.PrimitiveType] = func(
		yield func(int16, rs.PrimitiveType) bool,
	) {
		var ix2info iter.Seq2[int16, ColumnInfo] = maps.All(colmap)
		for ix, info := range ix2info {
			var typ rs.PrimitiveType = info.PrimitiveType
			yield(ix, typ)
		}
	}
	var ix2typMap map[int16]rs.PrimitiveType = maps.Collect(ix2typs)
	var ix2col2val func(
		int16, PgColumn,
	) IO[sw.Value] = ColumnIndexToColumnToValueNewDefaultFromMap(
		ix2typMap,
	)

	return PgRowsToValues(
		ctx,
		rows,
		ix2name,
		ix2col2val,
	)
}

func ColumnMapToPgRows(
	colmap map[int16]ColumnInfo,
) func(iter.Seq2[PgRow, error]) IO[iter.Seq2[map[string]sw.Value, error]] {
	return func(
		rows iter.Seq2[PgRow, error],
	) IO[iter.Seq2[map[string]sw.Value, error]] {
		return func(
			ctx context.Context,
		) (iter.Seq2[map[string]sw.Value, error], error) {
			return PgRowsToValuesFromColumnMap(
				ctx,
				rows,
				colmap,
			), nil
		}
	}
}

var MapToNameResolver func(map[int16]string) func(int16) (string, error) = rs.
	GetValueByKeyFromMap[int16, string](
	func(ix int16) error {
		return fmt.Errorf("%w: %v", ErrColumnNameNotFound, ix)
	},
)

func NamesToNameResolverMap(
	names iter.Seq[string],
) map[int16]string {
	var pairs iter.Seq2[int16, string] = func(
		yield func(int16, string) bool,
	) {
		var ix int16 = 0
		for name := range names {
			yield(ix, name)
			ix += 1
		}
	}
	return maps.Collect(pairs)
}

func JoinedNamesToNameResolverMap(
	split string,
	joined string,
) map[int16]string {
	var splited []string = strings.Split(joined, split)
	var i iter.Seq[string] = slices.Values(splited)
	return NamesToNameResolverMap(i)
}

var JoinedNamesToNameResolverMapDefault func(
	string,
) map[int16]string = Curry(JoinedNamesToNameResolverMap)(",")
