package types

import (
	"fmt"
	"strconv"
	"strings"
)

// StaticValue is an implementation of ValueLoader which is just
// a static float64 value. Using of this type allows to avoid extra
// function calls, sometimes.
type StaticValue float64

// Load implements ValueLoader.
//go:nosplit
func (r StaticValue) Load() float64 {
	return (float64)(r)
}

// FuncValue is just a function wrapper which implements ValueLoader.
type FuncValue func() float64

// Load implements ValueLoader.
func (r FuncValue) Load() float64 {
	return r()
}

// ValueLoader is something able to return a value of the variable.
type ValueLoader interface {
	// Load returns the value of the variable.
	Load() float64
}

// ParseValue returns a ValueLoader for variable or constant passed in `value`.
func ParseValue(value string, symResolver SymbolResolver) (ValueLoader, error) {
	switch {
	case strings.HasPrefix(value, "0x"):
		v, err := strconv.ParseInt(value[2:], 16, 64)
		if err == nil {
			return StaticValue(v), nil
		}
	case strings.HasPrefix(value, "h"):
		v, err := strconv.ParseInt(value[1:], 16, 64)
		if err == nil {
			return StaticValue(v), nil
		}
	case strings.HasPrefix(value, "b"):
		v, err := strconv.ParseInt(value[1:], 2, 64)
		if err == nil {
			return StaticValue(v), nil
		}
	case strings.HasPrefix(value, "o"):
		v, err := strconv.ParseInt(value[1:], 8, 64)
		if err == nil {
			return StaticValue(v), nil
		}
	default:
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return StaticValue(v), nil
		} else {
			v, err := strconv.ParseFloat(value, 64)
			if err == nil {
				return StaticValue(v), nil
			}
		}
	}

	if symResolver == nil {
		return nil, fmt.Errorf("valueLoader is not set, and cannot valueLoader symbol: '%s'", value)
	}
	valueLoader, err := symResolver.Resolve(value)
	if err != nil {
		return nil, fmt.Errorf("unable to valueLoader symbol '%s': %w", value, err)
	}
	return valueLoader, nil
}
