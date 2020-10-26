package rpn

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/xaionaro-go/rpn/types"
)

var (
	_ types.Expr = &Expr{}
)

// Expr is an implementation of types.Expr which uses LLVM JIT code to
// evaluate the expression.
//
// WARNING! This is unsafe implementation, do not use it if you haven't
// checked Ops.Compile by yourself!
type Expr struct {
	Description           string
	Code                  func() float64
	Syms                  []Symbol
	ResultCache           types.NullFloat64
	IsMemoizationEnabled  bool
	stack                 []float64
	values                []float64
	nonStaticValueIndices []int
}

// Symbol provides information how to extract the value and what name
// the symbol (of the expression) has. Symbol -- is anything except for
// operations signs.
type Symbol struct {
	StaticValue types.NullFloat64
	ValueLoader types.ValueLoader
	Name        string
}

// Load returns the current value of the symbol
func (sym *Symbol) Load() float64 {
	if sym.StaticValue.Valid {
		return sym.StaticValue.Float64
	}
	return sym.ValueLoader.Load()
}

// Eval implements types.Expr
func (expr *Expr) Eval() float64 {
	if !expr.IsMemoizationEnabled {
		return expr.eval()
	}

	if expr.ResultCache.Valid {
		return expr.ResultCache.Float64
	}

	r := expr.eval()

	expr.ResultCache.Float64 = r
	expr.ResultCache.Valid = true

	return r
}

func (expr *Expr) eval() float64 {
	values := expr.values
	for _, idx := range expr.nonStaticValueIndices {
		values[idx] = expr.Syms[idx].Load()
	}
	return expr.Code()
}

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
//
// input example: "z x y + *"
// calculation interpretation: z * (x + y)
func Parse(expression string, symResolver types.SymbolResolver) (*Expr, error) {

	ops := Ops{}
	expr := &Expr{
		Description: expression,
	}
	parts := strings.Split(expression, " ")
	for _, part := range parts {
		if part == "" {
			continue
		}
		op := types.ParseOp(part)

		if op != types.OpUndefined {
			ops = append(ops, op)
			continue
		}
		ops = append(ops, types.OpFetch)

		valueLoader, err := types.ParseValue(part, symResolver)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value '%s': %w", part, err)
		}

		sym := Symbol{
			ValueLoader: valueLoader,
			Name:        part,
		}
		if f, ok := valueLoader.(types.StaticValue); ok {
			sym.StaticValue.Valid = true
			sym.StaticValue.Float64 = f.Load()
		} else {
			expr.nonStaticValueIndices = append(expr.nonStaticValueIndices, len(expr.Syms))
		}
		expr.Syms = append(expr.Syms, sym)
	}

	expr.stack = make([]float64, len(expr.Syms))
	expr.values = make([]float64, len(expr.Syms))

	for idx, sym := range expr.Syms {
		if sym.StaticValue.Valid {
			expr.values[idx] = sym.StaticValue.Float64
		}
	}

	var cleanup func()
	expr.Code, cleanup = ops.Compile(expr.stack, expr.values)
	runtime.SetFinalizer(expr, func(expr *Expr) {
		cleanup()
	})
	return expr, nil
}

// String implements types.Expr
func (expr *Expr) String() string {
	return expr.Description
}

// EnableMemoization implements types.Expr
func (expr *Expr) EnableMemoization(newValue bool) (oldValue bool) {
	oldValue = expr.IsMemoizationEnabled
	expr.IsMemoizationEnabled = newValue
	return
}
