package writer

// This file is generated using a2vsw.tmpl. NEVER EDIT.

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func AnyToVal(a any) Value {
	switch t := a.(type) {

	case string:
		return StringToValue(t)

	case sql.Null[string]:
		return NullableStringToValue(t)
	case *string:
		var conv func(sql.Null[string]) Value = NullableStringToValue
		var nval sql.Null[string] = ConvertToNullable(t)
		return conv(nval)

	case []byte:
		return BytesToValue(t)

	case sql.Null[[]byte]:
		return NullableBytesToValue(t)
	case *[]byte:
		var conv func(sql.Null[[]byte]) Value = NullableBytesToValue
		var nval sql.Null[[]byte] = ConvertToNullable(t)
		return conv(nval)

	case int32:
		return IntToValue(t)

	case sql.Null[int32]:
		return NullableIntToValue(t)
	case *int32:
		var conv func(sql.Null[int32]) Value = NullableIntToValue
		var nval sql.Null[int32] = ConvertToNullable(t)
		return conv(nval)

	case int64:
		return LongToValue(t)

	case sql.Null[int64]:
		return NullableLongToValue(t)
	case *int64:
		var conv func(sql.Null[int64]) Value = NullableLongToValue
		var nval sql.Null[int64] = ConvertToNullable(t)
		return conv(nval)

	case float32:
		return FloatToValue(t)

	case sql.Null[float32]:
		return NullableFloatToValue(t)
	case *float32:
		var conv func(sql.Null[float32]) Value = NullableFloatToValue
		var nval sql.Null[float32] = ConvertToNullable(t)
		return conv(nval)

	case float64:
		return DoubleToValue(t)

	case sql.Null[float64]:
		return NullableDoubleToValue(t)
	case *float64:
		var conv func(sql.Null[float64]) Value = NullableDoubleToValue
		var nval sql.Null[float64] = ConvertToNullable(t)
		return conv(nval)

	case bool:
		return BooleanToValue(t)

	case sql.Null[bool]:
		return NullableBooleanToValue(t)
	case *bool:
		var conv func(sql.Null[bool]) Value = NullableBooleanToValue
		var nval sql.Null[bool] = ConvertToNullable(t)
		return conv(nval)

	case struct{}:
		return NullToValue(t)

	case time.Time:
		return TimeToValue(t)

	case sql.Null[time.Time]:
		return NullableTimeToValue(t)
	case *time.Time:
		var conv func(sql.Null[time.Time]) Value = NullableTimeToValue
		var nval sql.Null[time.Time] = ConvertToNullable(t)
		return conv(nval)

	case uuid.UUID:
		return UuidToValue(t)

	case sql.Null[uuid.UUID]:
		return NullableUuidToValue(t)
	case *uuid.UUID:
		var conv func(sql.Null[uuid.UUID]) Value = NullableUuidToValue
		var nval sql.Null[uuid.UUID] = ConvertToNullable(t)
		return conv(nval)

	case nil:
		return NullToValue(struct{}{})
	default:
	}
	return InvalidValueFromErr(fmt.Errorf("%w: %v", ErrInvalidValue, a))
}
