package rpn

import (
	"fmt"
	"math"
	"strings"

	"github.com/xaionaro-go/rpn/internal"
	"github.com/xaionaro-go/rpn/types"
)

var (
	_ types.Expr = &Expr{}
)

// x y z + +
// y + z
// y+z + x

// x y + z +
// x + y
// x+y + z

// Expr is an implementation of types.Expr which tries to present the
// expression in a flat format (as a slice) to avoid extra performance
// penalties on traversing a tree (in comparison to the "calltree"
// implementation)
type Expr struct {
	CallNodes            []func()
	ResultCache          types.NullFloat64
	RAM                  []float64
	IsMemoizationEnabled bool
	Description          string
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
	for _, callNode := range expr.CallNodes {
		callNode()
	}
	return expr.RAM[len(expr.RAM)-1]
}

type value struct {
	internal.ParsedValue
	RAMIdx int
}

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
//
// input example: "z x y + *"
// calculation interpretation: z * (x + y)
func Parse(expression string, symResolver types.SymbolResolver) (*Expr, error) {
	expr := &Expr{
		Description: expression,
	}
	parts := strings.Split(expression, " ")
	values := make([]value, 0, 2)
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

			values = append(values, value{ParsedValue: parsedValue, RAMIdx: -1})
			continue
		}

		unusedSymsCount := len(values)
		if unusedSymsCount < 2 {
			return nil, fmt.Errorf("expected at least 2 values in stack, but found only %d (partIdx: %d; expression: '%s')", unusedSymsCount, partIdx, expression)
		}

		lhsSym := values[len(values)-2]
		rhsSym := values[len(values)-1]
		values = values[:len(values)-2]

		ramIdx := len(expr.RAM)
		expr.RAM = append(expr.RAM, float64(0))
		values = append(values, value{RAMIdx: ramIdx})

		switch {
		case lhsSym.ConstValue.Valid && rhsSym.ConstValue.Valid:
			lhs, rhs := lhsSym.ConstValue.Float64, rhsSym.ConstValue.Float64
			expr.RAM[ramIdx] = op.Eval(lhs, rhs)

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 + rhsSym.FuncValue()
			})
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 + expr.RAM[rhsSym.RAMIdx]
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() + rhsSym.ConstValue.Float64
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] + rhsSym.ConstValue.Float64
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() + rhsSym.FuncValue()
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() + expr.RAM[rhsSym.RAMIdx]
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] + rhsSym.FuncValue()
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpPlus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] + expr.RAM[rhsSym.RAMIdx]
			})

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 - rhsSym.FuncValue()
			})
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 - expr.RAM[rhsSym.RAMIdx]
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() - rhsSym.ConstValue.Float64
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] - rhsSym.ConstValue.Float64
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() - rhsSym.FuncValue()
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() - expr.RAM[rhsSym.RAMIdx]
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] - rhsSym.FuncValue()
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpMinus:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] - expr.RAM[rhsSym.RAMIdx]
			})

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 * rhsSym.FuncValue()
			})
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 * expr.RAM[rhsSym.RAMIdx]
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() * rhsSym.ConstValue.Float64
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] * rhsSym.ConstValue.Float64
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() * rhsSym.FuncValue()
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() * expr.RAM[rhsSym.RAMIdx]
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] * rhsSym.FuncValue()
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpMultiply:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] * expr.RAM[rhsSym.RAMIdx]
			})

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 / rhsSym.FuncValue()
			})
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.ConstValue.Float64 / expr.RAM[rhsSym.RAMIdx]
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() / rhsSym.ConstValue.Float64
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] / rhsSym.ConstValue.Float64
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() / rhsSym.FuncValue()
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = lhsSym.FuncValue() / expr.RAM[rhsSym.RAMIdx]
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] / rhsSym.FuncValue()
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpDivide:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = expr.RAM[lhsSym.RAMIdx] / expr.RAM[rhsSym.RAMIdx]
			})

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(lhsSym.ConstValue.Float64, rhsSym.FuncValue())
			})
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(lhsSym.ConstValue.Float64, expr.RAM[rhsSym.RAMIdx])
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(lhsSym.FuncValue(), rhsSym.ConstValue.Float64)
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(expr.RAM[lhsSym.RAMIdx], rhsSym.ConstValue.Float64)
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(lhsSym.FuncValue(), rhsSym.FuncValue())
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(lhsSym.FuncValue(), expr.RAM[rhsSym.RAMIdx])
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(expr.RAM[lhsSym.RAMIdx], rhsSym.FuncValue())
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpPower:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[ramIdx] = math.Pow(expr.RAM[lhsSym.RAMIdx], expr.RAM[rhsSym.RAMIdx])
			})

		case lhsSym.ConstValue.Valid && rhsSym.FuncValue != nil && op == types.OpIf:
			if lhsSym.ConstValue.Float64 > 0 {
				expr.CallNodes = append(expr.CallNodes, func() {
					expr.RAM[ramIdx] = rhsSym.FuncValue()
				})
			}
		case lhsSym.ConstValue.Valid && rhsSym.FuncValue == nil && op == types.OpIf:
			if lhsSym.ConstValue.Float64 > 0 {
				expr.CallNodes = append(expr.CallNodes, func() {
					expr.RAM[ramIdx] = expr.RAM[rhsSym.RAMIdx]
				})
			}
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue != nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if lhsSym.FuncValue() > 0 {
					expr.RAM[ramIdx] = rhsSym.ConstValue.Float64
					return
				}
				expr.RAM[ramIdx] = 0
			})
		case rhsSym.ConstValue.Valid && lhsSym.FuncValue == nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if expr.RAM[lhsSym.RAMIdx] > 0 {
					expr.RAM[ramIdx] = rhsSym.ConstValue.Float64
					return
				}
				expr.RAM[ramIdx] = 0
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue != nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if lhsSym.FuncValue() > 0 {
					expr.RAM[ramIdx] = rhsSym.FuncValue()
					return
				}
				expr.RAM[ramIdx] = 0
			})
		case lhsSym.FuncValue != nil && rhsSym.FuncValue == nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if lhsSym.FuncValue() > 0 {
					expr.RAM[ramIdx] = expr.RAM[rhsSym.RAMIdx]
					return
				}
				expr.RAM[ramIdx] = 0
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue != nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if expr.RAM[lhsSym.RAMIdx] > 0 {
					expr.RAM[ramIdx] = rhsSym.FuncValue()
					return
				}
				expr.RAM[ramIdx] = 0
			})
		case lhsSym.FuncValue == nil && rhsSym.FuncValue == nil && op == types.OpIf:
			expr.CallNodes = append(expr.CallNodes, func() {
				if expr.RAM[lhsSym.RAMIdx] > 0 {
					expr.RAM[ramIdx] = expr.RAM[rhsSym.RAMIdx]
					return
				}
				expr.RAM[ramIdx] = 0
			})
		}
	}

	if len(values) == 1 && len(expr.RAM) == 0 {
		// This is the case when no operators is given but just a value only
		expr.RAM = append(expr.RAM, float64(0))
		value := values[0]
		switch {
		case value.ConstValue.Valid:
			expr.RAM[0] = value.ConstValue.Float64
		case value.FuncValue != nil:
			expr.CallNodes = append(expr.CallNodes, func() {
				expr.RAM[0] = value.FuncValue()
			})
		default:
			panic("should not happen")
		}
	}

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
