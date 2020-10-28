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

// Expr is an implementation of types.Expr which tries to precalculate as
// much as possible and the rest is stored as a function which directly
// calls another functions and calculates the results.
type Expr struct {
	Description          string
	IsMemoizationEnabled bool
	RootFunc             func() float64
	ResultCache          types.NullFloat64
}

// Eval implements types.Expr
func (expr *Expr) Eval() float64 {
	if !expr.IsMemoizationEnabled {
		return expr.RootFunc()
	}

	if expr.ResultCache.Valid {
		return expr.ResultCache.Float64
	}

	r := expr.RootFunc()
	expr.ResultCache.Float64 = r
	expr.ResultCache.Valid = true

	return r
}

type stack []*internal.ParsedValue

func (s *stack) Push(node internal.ParsedValue) *internal.ParsedValue {
	*s = append(*s, &node)
	return s.First()
}

func (s *stack) First() *internal.ParsedValue {
	return (*s)[len(*s)-1]
}

func (s *stack) Pop() *internal.ParsedValue {
	r := s.First()
	*s = (*s)[:len(*s)-1]
	return r
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
	stack := stack{}
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
			stack.Push(parsedValue)
			continue
		}

		if len(stack) < 2 {
			return nil, fmt.Errorf("invalid expression '%s' at part index %d: expected at least two entries in the stack", expression, partIdx)
		}
		rhs := stack.Pop()
		lhs := stack.Pop()

		switch op {
		case types.OpPlus:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: lhs.ConstValue.Float64 + rhs.ConstValue.Float64,
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() + rhs.ConstValue.Float64
					},
				})
			case lhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.ConstValue.Float64 + rhs.FuncValue()
					},
				})
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() + rhs.FuncValue()
					},
				})
			}
		case types.OpMinus:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: lhs.ConstValue.Float64 - rhs.ConstValue.Float64,
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() - rhs.ConstValue.Float64
					},
				})
			case lhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.ConstValue.Float64 - rhs.FuncValue()
					},
				})
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() - rhs.FuncValue()
					},
				})
			}
		case types.OpMultiply:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: lhs.ConstValue.Float64 * rhs.ConstValue.Float64,
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() * rhs.ConstValue.Float64
					},
				})
			case lhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.ConstValue.Float64 * rhs.FuncValue()
					},
				})
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() * rhs.FuncValue()
					},
				})
			}
		case types.OpDivide:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: lhs.ConstValue.Float64 / rhs.ConstValue.Float64,
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() / rhs.ConstValue.Float64
					},
				})
			case lhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.ConstValue.Float64 / rhs.FuncValue()
					},
				})
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return lhs.FuncValue() / rhs.FuncValue()
					},
				})
			}
		case types.OpPower:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: math.Pow(lhs.ConstValue.Float64, rhs.ConstValue.Float64),
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return math.Pow(lhs.FuncValue(), rhs.ConstValue.Float64)
					},
				})
			case lhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return math.Pow(lhs.ConstValue.Float64, rhs.FuncValue())
					},
				})
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						return math.Pow(lhs.FuncValue(), rhs.FuncValue())
					},
				})
			}
		case types.OpIf:
			switch {
			case lhs.ConstValue.Valid && rhs.ConstValue.Valid:
				v := float64(0)
				if lhs.ConstValue.Float64 > 0 {
					v = rhs.ConstValue.Float64
				}
				stack.Push(internal.ParsedValue{
					ConstValue: types.NullFloat64{
						Valid:   true,
						Float64: v,
					},
				})
			case rhs.ConstValue.Valid:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						if lhs.FuncValue() > 0 {
							return rhs.ConstValue.Float64
						}
						return 0
					},
				})
			case lhs.ConstValue.Valid:
				if lhs.ConstValue.Float64 > 0 {
					stack.Push(internal.ParsedValue{
						FuncValue: func() float64 {
							return rhs.FuncValue()
						},
					})
				} else {
					stack.Push(internal.ParsedValue{
						ConstValue: types.NullFloat64{
							Float64: 0,
							Valid:   true,
						},
					})
				}
			default:
				stack.Push(internal.ParsedValue{
					FuncValue: func() float64 {
						if lhs.FuncValue() > 0 {
							return rhs.FuncValue()
						}
						return 0
					},
				})
			}
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("expected stack length is 1, but got %d", len(stack))
	}
	rootCallNode := stack[0]

	if rootCallNode.ConstValue.Valid {
		expr.RootFunc = func() float64 {
			return rootCallNode.ConstValue.Float64
		}
	} else {
		expr.RootFunc = rootCallNode.FuncValue
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
