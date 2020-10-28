package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xaionaro-go/rpn/types"
)

// ParsedValue is for
type ParsedValue struct {
	ConstValue types.NullFloat64
	FuncValue  types.FuncValue
}

// Load implements ValueLoader
func (v *ParsedValue) Load() float64 {
	if v.ConstValue.Valid {
		return v.ConstValue.Float64
	}
	return v.FuncValue()
}

// ParseValue returns a ValueLoader for variable or constant passed in `value`.
func ParseValue(value string, symResolver types.SymbolResolver) (ParsedValue, error) {
	var (
		v   float64
		i   int64
		err error
	)
	switch {
	case strings.HasPrefix(value, "0x"):
		i, err = strconv.ParseInt(value[2:], 16, 64)
		v = float64(i)
	case strings.HasPrefix(value, "h"):
		i, err = strconv.ParseInt(value[1:], 16, 64)
		v = float64(i)
	case strings.HasPrefix(value, "b"):
		i, err = strconv.ParseInt(value[1:], 2, 64)
		v = float64(i)
	case strings.HasPrefix(value, "o"):
		i, err = strconv.ParseInt(value[1:], 8, 64)
		v = float64(i)
	default:
		v, err = strconv.ParseFloat(value, 64)
	}
	if err == nil {
		return ParsedValue{
			ConstValue: types.NullFloat64{
				Float64: v,
				Valid:   true,
			},
		}, nil
	}

	if symResolver == nil {
		return ParsedValue{}, fmt.Errorf("valueLoader is not set, and cannot valueLoader symbol: '%s'", value)
	}
	valueLoader, err := symResolver.Resolve(value)
	if err != nil {
		return ParsedValue{}, fmt.Errorf("unable to valueLoader symbol '%s': %w", value, err)
	}

	r := ParsedValue{}
	switch valueLoader := valueLoader.(type) {
	case types.StaticValue:
		r.ConstValue = types.NullFloat64{
			Float64: float64(valueLoader),
			Valid:   true,
		}
	case types.FuncValue:
		r.FuncValue = valueLoader
	default:
		r.FuncValue = valueLoader.Load
	}
	return r, nil
}
