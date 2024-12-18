package pgcopy2vals

// This file is generated using a2vsw.tmpl. NEVER EDIT.

import (
	"strings"

	rs "github.com/takanoriyanagitani/go-rowlike2sqlvalues"
	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

var type2convGenMap TypeToConverterMap = map[rs.PrimitiveType]ConfigToConverter{

	rs.PrimitiveBytes: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueBytes() })
		}
	},

	rs.NullBytes: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullBytes() })
		}
	},

	rs.PrimitiveInt: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueInt() })
		}
	},

	rs.NullInt: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullInt() })
		}
	},

	rs.PrimitiveLong: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueLong() })
		}
	},

	rs.NullLong: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullLong() })
		}
	},

	rs.PrimitiveFloat: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueFloat() })
		}
	},

	rs.NullFloat: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullFloat() })
		}
	},

	rs.PrimitiveDouble: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueDouble() })
		}
	},

	rs.NullDouble: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullDouble() })
		}
	},

	rs.PrimitiveBoolean: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueBoolean() })
		}
	},

	rs.NullBoolean: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullBoolean() })
		}
	},

	rs.PrimitiveTime: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueTime() })
		}
	},

	rs.NullTime: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullTime() })
		}
	},

	rs.PrimitiveUuid: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueUuid() })
		}
	},

	rs.NullUuid: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNullUuid() })
		}
	},

	rs.PrimitiveString: func(
		cfg PgColumnConverterConfig,
	) PgColumnToValue {
		var buf strings.Builder
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value {
				buf.Reset()
				var chk func(string) error = cfg.StringChecker
				return c.ToValueString(chk, &buf)
			})
		}
	},

	rs.NullString: func(
		cfg PgColumnConverterConfig,
	) PgColumnToValue {
		var buf strings.Builder
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value {
				buf.Reset()
				var chk func(string) error = cfg.StringChecker
				return c.ToValueNullString(chk, &buf)
			})
		}
	},

	rs.PrimitiveNull: func(
		_ PgColumnConverterConfig,
	) PgColumnToValue {
		return func(c PgColumn) IO[sw.Value] {
			return OfFn(func() sw.Value { return c.ToValueNull() })
		}
	},
}
