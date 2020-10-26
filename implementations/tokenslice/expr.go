package rpn

import (
	"fmt"
	"math"
	"strings"

	"github.com/xaionaro-go/rpn/types"
)

var (
	_ types.Expr = &Expr{}
)

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
	symIdx := len(expr.Syms) - 1
	stackLen := 0
	syms := expr.Syms
	stack := expr.evalStack
	ops := expr.Ops
	for _, op := range ops {
		var lhs, rhs float64
		switch stackLen {
		case 0:
			rhs = syms[symIdx].Load()
			lhs = syms[symIdx-1].Load()
			symIdx -= 2
		case 1:
			rhs = syms[symIdx].Load()
			symIdx--
			stackLen--
			lhs = stack[stackLen]
		default:
			stackLen--
			rhs = stack[stackLen]
			stackLen--
			lhs = stack[stackLen]
		}

		r := float64(0)
		switch op {
		case types.OpPlus:
			r = lhs + rhs
		case types.OpMinus:
			r = lhs - rhs
		case types.OpMultiply:
			r = lhs * rhs
		case types.OpDivide:
			r = lhs / rhs
		case types.OpPower:
			r = math.Pow(lhs, rhs)
		case types.OpIf:
			if lhs > 0 {
				r = rhs
			}
		default:
			panic("should not happened")
		}

		if symIdx < 0 {
			return r
		}

		stack[stackLen] = r
		stackLen++
	}
	if stackLen == 0 {
		return expr.Syms[0].StaticValue.Float64
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
	for _, part := range parts {
		if part == "" {
			continue
		}
		op := types.ParseOp(part)

		if op == types.OpUndefined {
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
			}
			expr.Syms = append(expr.Syms, sym)
			continue
		}

		lhsSym := expr.Syms[len(expr.Syms)-2]
		rhsSym := expr.Syms[len(expr.Syms)-1]
		if !lhsSym.StaticValue.Valid || !rhsSym.StaticValue.Valid {
			expr.Ops = append(expr.Ops, op)
			continue
		}
		lhs, rhs := lhsSym.StaticValue.Float64, rhsSym.StaticValue.Float64
		var r float64
		switch op {
		case types.OpPlus:
			r = lhs + rhs
		case types.OpMinus:
			r = lhs - rhs
		case types.OpMultiply:
			r = lhs * rhs
		case types.OpDivide:
			r = lhs / rhs
		case types.OpPower:
			r = math.Pow(lhs, rhs)
		case types.OpIf:
			if lhs > 0 {
				r = rhs
			}
		default:
			panic("should not happened")
		}
		expr.Syms[len(expr.Syms)-2] = Symbol{
			StaticValue: types.NullFloat64{
				Valid:   true,
				Float64: r,
			},
			Name: "( " + lhsSym.Name + " " + op.String() + " " + rhsSym.Name + " )",
		}
		expr.Syms = expr.Syms[:len(expr.Syms)-1]
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
