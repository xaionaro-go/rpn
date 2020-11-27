package rpn

import (
	"fmt"
	"strings"

	"github.com/xaionaro-go/rpn/internal"
	"github.com/xaionaro-go/rpn/types"
)

var (
	_ types.Expr = &Expr{}
)

func init() {
	//panic("This package should not be used, it works wrong")
}

// Expr is an implementation of types.Expr which tries to present the
// expression in a flat format (as a slice) to avoid extra performance
// penalties on traversing a tree (in comparison to the "exprtree"
// implementation)
type Expr struct {
	Ops                  []types.Op
	Syms                 []Symbol
	ResultCache          types.NullFloat64
	IsMemoizationEnabled bool
	evalStack            []float64
}

// Symbol provides information how to extract the value and what name
// the symbol (of the expression) has. Symbol -- is anything except for
// operations signs.
type Symbol struct {
	internal.ParsedValue
	Name string
}

// Eval implements types.Expr
func (expr *Expr) Eval() float64 {
	if !expr.IsMemoizationEnabled {
		r := expr.eval()
		return r
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
	symIdx := 0
	stackLen := 0
	syms := expr.Syms
	stack := expr.evalStack
	ops := expr.Ops
	for _, op := range ops {
		if op == types.OpFetch {
			stack[stackLen] = syms[symIdx].Load()
			symIdx++
			stackLen++
			continue
		}

		stackLen--
		rhs := stack[stackLen]
		stackLen--
		lhs := stack[stackLen]

		r := op.Eval(lhs, rhs)

		if symIdx < 0 {
			return r
		}

		stack[stackLen] = r
		stackLen++
	}
	if stackLen == 0 {
		return expr.Syms[0].ConstValue.Float64
	}

	return stack[0]
}

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
//
// input example: "z x y + *"
// calculation interpretation: z * (x + y)
func Parse(expression string, symResolver types.SymbolResolver) (*Expr, error) {
	expr := &Expr{}
	parts := strings.Split(expression, " ")
	for partIdx, part := range parts {
		if part == "" {
			continue
		}
		op := types.ParseOp(part)

		if op == types.OpUndefined {
			parsedValue, err := internal.ParseValue(part, symResolver)
			if err != nil {
				return nil, fmt.Errorf("unable to parse value '%s': %w", part, err)
			}

			sym := Symbol{
				Name:        part,
				ParsedValue: parsedValue,
			}
			expr.Syms = append(expr.Syms, sym)
			expr.Ops = append(expr.Ops, types.OpFetch)
			continue
		}

		unusedSymsCount := len(expr.Syms)*2 - len(expr.Ops)
		if unusedSymsCount < 2 {
			return nil, fmt.Errorf("expected at least 2 values in stack, but found only %d (partIdx: %d; expression: '%s')", unusedSymsCount, partIdx, expression)
		}

		expr.Ops = append(expr.Ops, op)
	}
	expr.evalStack = make([]float64, len(expr.Syms))
	return expr, nil
}

// String implements types.Expr
func (expr *Expr) String() string {
	return fmt.Sprintf("%v with %v", expr.Ops, expr.Syms)
}

// EnableMemoization implements types.Expr
func (expr *Expr) EnableMemoization(newValue bool) (oldValue bool) {
	oldValue = expr.IsMemoizationEnabled
	expr.IsMemoizationEnabled = newValue
	return
}
