package rowlike2sqlvals

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

var (
	ErrInvalidType error = errors.New("invalid type")
)

type PrimitiveType string

func GetValueByKeyFromMap[K comparable, V any](
	onMissing func(K) error,
) func(map[K]V) func(K) (V, error) {
	return func(m map[K]V) func(K) (V, error) {
		return func(key K) (v V, e error) {
			val, found := m[key]
			switch found {
			case true:
				return val, nil
			default:
				return v, onMissing(key)
			}
		}
	}
}

const (
	PrimitiveUnknown PrimitiveType = "UNKNOWN"

	PrimitiveString  PrimitiveType = "string"
	PrimitiveBytes   PrimitiveType = "bytes"
	PrimitiveInt     PrimitiveType = "int"
	PrimitiveLong    PrimitiveType = "long"
	PrimitiveFloat   PrimitiveType = "float"
	PrimitiveDouble  PrimitiveType = "double"
	PrimitiveBoolean PrimitiveType = "boolean"
	PrimitiveNull    PrimitiveType = "null"

	PrimitiveTime PrimitiveType = "time"
	PrimitiveUuid PrimitiveType = "uuid"
)

const (
	NullString  PrimitiveType = "null-string"
	NullBytes   PrimitiveType = "null-bytes"
	NullInt     PrimitiveType = "null-int"
	NullLong    PrimitiveType = "null-long"
	NullFloat   PrimitiveType = "null-float"
	NullDouble  PrimitiveType = "null-double"
	NullBoolean PrimitiveType = "null-boolean"

	NullTime PrimitiveType = "null-time"
	NullUuid PrimitiveType = "null-uuid"
)

var primitiveTypes []PrimitiveType = slices.Collect(
	func(yield func(PrimitiveType) bool) {

		yield(PrimitiveUnknown)
		yield(PrimitiveString)
		yield(PrimitiveBytes)
		yield(PrimitiveInt)
		yield(PrimitiveLong)
		yield(PrimitiveFloat)
		yield(PrimitiveDouble)
		yield(PrimitiveBoolean)
		yield(PrimitiveNull)

		yield(NullString)
		yield(NullBytes)
		yield(NullInt)
		yield(NullLong)
		yield(NullFloat)
		yield(NullDouble)
		yield(NullBoolean)

		yield(PrimitiveTime)
		yield(NullTime)

		yield(PrimitiveUuid)
		yield(NullUuid)

	},
)

var typ2stringMap map[PrimitiveType]string = maps.Collect(
	func(yield func(PrimitiveType, string) bool) {
		for _, typ := range primitiveTypes {
			var s string = string(typ)
			yield(typ, s)
		}
	},
)

var string2typMap map[string]PrimitiveType = maps.Collect(
	func(yield func(string, PrimitiveType) bool) {
		for typ, s := range typ2stringMap {
			yield(s, typ)
		}
	},
)

type StringToPrimitiveType func(string) (PrimitiveType, error)

var StringToType StringToPrimitiveType = GetValueByKeyFromMap[string, PrimitiveType](
	func(key string) error {
		return fmt.Errorf("%w: %s", ErrInvalidType, key)
	},
)(string2typMap)
