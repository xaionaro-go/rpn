package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/rpn/types"
)

// DummyResolver is just a variables resolver for tests.
type DummyResolver struct {
	T *testing.T
}

// Resolve implements types.SymbolResolver
func (r DummyResolver) Resolve(sym string) (types.ValueLoader, error) {
	switch sym {
	case "x0":
		return types.FuncValue(func() float64 {
			return 2
		}), nil
	case "x1":
		return types.FuncValue(func() float64 {
			return 3
		}), nil
	case "y":
		return types.StaticValue(4), nil
	case "z":
		return types.FuncValue(func() float64 {
			return 1
		}), nil
	}
	require.FailNow(r.T, fmt.Sprintf("should not happen: '%s'", sym))
	return nil, nil
}
