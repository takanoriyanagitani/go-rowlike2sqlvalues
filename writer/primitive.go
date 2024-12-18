package writer

import (
	"database/sql"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var (
	ErrUnimplementedWriter error = errors.New("unimplemented writer")
)

type PrimitiveWriter struct {
	StringWriter  func(string) IO[Void]
	BytesWriter   func([]byte) IO[Void]
	IntWriter     func(int32) IO[Void]
	LongWriter    func(int64) IO[Void]
	FloatWriter   func(float32) IO[Void]
	DoubleWriter  func(float64) IO[Void]
	BooleanWriter func(bool) IO[Void]
	NullWriter    func(struct{}) IO[Void]

	TimeWriter func(time.Time) IO[Void]
	UuidWriter func(uuid.UUID) IO[Void]
}

var InvalidWtrRes IO[Void] = Err[Void](ErrUnimplementedWriter)

var PrimitiveWriterDefault PrimitiveWriter = PrimitiveWriter{
	StringWriter:  func(_ string) IO[Void] { return InvalidWtrRes },
	BytesWriter:   func(_ []byte) IO[Void] { return InvalidWtrRes },
	IntWriter:     func(_ int32) IO[Void] { return InvalidWtrRes },
	LongWriter:    func(_ int64) IO[Void] { return InvalidWtrRes },
	FloatWriter:   func(_ float32) IO[Void] { return InvalidWtrRes },
	DoubleWriter:  func(_ float64) IO[Void] { return InvalidWtrRes },
	BooleanWriter: func(_ bool) IO[Void] { return InvalidWtrRes },
	NullWriter:    func(_ struct{}) IO[Void] { return InvalidWtrRes },

	TimeWriter: func(_ time.Time) IO[Void] { return InvalidWtrRes },
	UuidWriter: func(_ uuid.UUID) IO[Void] { return InvalidWtrRes },
}

type NullableWriter struct {
	StringWriter  func(sql.Null[string]) IO[Void]
	BytesWriter   func(sql.Null[[]byte]) IO[Void]
	IntWriter     func(sql.Null[int32]) IO[Void]
	LongWriter    func(sql.Null[int64]) IO[Void]
	FloatWriter   func(sql.Null[float32]) IO[Void]
	DoubleWriter  func(sql.Null[float64]) IO[Void]
	BooleanWriter func(sql.Null[bool]) IO[Void]

	TimeWriter func(sql.Null[time.Time]) IO[Void]
	UuidWriter func(sql.Null[uuid.UUID]) IO[Void]
}

var NullableWriterDefault NullableWriter = NullableWriter{
	StringWriter:  func(_ sql.Null[string]) IO[Void] { return InvalidWtrRes },
	BytesWriter:   func(_ sql.Null[[]byte]) IO[Void] { return InvalidWtrRes },
	IntWriter:     func(_ sql.Null[int32]) IO[Void] { return InvalidWtrRes },
	LongWriter:    func(_ sql.Null[int64]) IO[Void] { return InvalidWtrRes },
	FloatWriter:   func(_ sql.Null[float32]) IO[Void] { return InvalidWtrRes },
	DoubleWriter:  func(_ sql.Null[float64]) IO[Void] { return InvalidWtrRes },
	BooleanWriter: func(_ sql.Null[bool]) IO[Void] { return InvalidWtrRes },

	TimeWriter: func(_ sql.Null[time.Time]) IO[Void] { return InvalidWtrRes },
	UuidWriter: func(_ sql.Null[uuid.UUID]) IO[Void] { return InvalidWtrRes },
}

type ValueWriter struct {
	PrimitiveWriter
	NullableWriter
}

var ValueWriterDefault ValueWriter = ValueWriter{
	PrimitiveWriter: PrimitiveWriterDefault,
	NullableWriter:  NullableWriterDefault,
}

//go:generate go run internal/gen/any2val/main.go String string
//go:generate go run internal/gen/any2val/main.go Bytes []byte
//go:generate go run internal/gen/any2val/main.go Int int32
//go:generate go run internal/gen/any2val/main.go Long int64
//go:generate go run internal/gen/any2val/main.go Float float32
//go:generate go run internal/gen/any2val/main.go Double float64
//go:generate go run internal/gen/any2val/main.go Boolean bool
//go:generate go run internal/gen/any2val/main.go Null struct{}
//go:generate go run internal/gen/any2val/main.go Time time.Time
//go:generate go run internal/gen/any2val/main.go Uuid uuid.UUID
//go:generate gofmt -s -w .
type Value func(ValueWriter) IO[Void]

func InvalidValueFromErr(err error) Value {
	return func(_ ValueWriter) IO[Void] {
		return Err[Void](err)
	}
}

func PrimitiveWriterNew(w io.Writer) PrimitiveWriter {
	return PrimitiveWriter{
		StringWriter:  StringWriterNew(w),
		BytesWriter:   BytesWriterNew(w),
		IntWriter:     IntWriterNew(w),
		LongWriter:    LongWriterNew(w),
		FloatWriter:   FloatWriterNewDefault(w),
		DoubleWriter:  DoubleWriterNewDefault(w),
		BooleanWriter: BooleanWriterNew(w),
		NullWriter:    NullWriterNew(w),
		TimeWriter:    TimeWriterNew(w),
		UuidWriter:    UuidWriterNew(w),
	}
}

func NullableWriterNew(w io.Writer) NullableWriter {
	return NullableWriter{
		StringWriter:  NullStringWriterNew(w),
		BytesWriter:   NullBytesWriterNew(w),
		IntWriter:     NullIntWriterNew(w),
		LongWriter:    NullLongWriterNew(w),
		FloatWriter:   NullFloatWriterNewDefault(w),
		DoubleWriter:  NullDoubleWriterNewDefault(w),
		BooleanWriter: NullBooleanWriterNew(w),
		TimeWriter:    NullTimeWriterNew(w),
		UuidWriter:    NullUuidWriterNew(w),
	}
}

func ValueWriterNew(w io.Writer) ValueWriter {
	return ValueWriter{
		PrimitiveWriter: PrimitiveWriterNew(w),
		NullableWriter:  NullableWriterNew(w),
	}
}
