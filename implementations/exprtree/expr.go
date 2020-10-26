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

// Expr is an implementation of types.Expr which works are a tree of
// expressions, so you may access any sub-tree if you need it. This is
// the slowest implementation from this collection.
type Expr struct {
	LHS           *Expr
	RHS           *Expr
	Symbol        string
	ValueLoader   types.ValueLoader
	StaticValue   types.NullFloat64
	ResultCache   types.NullFloat64
	IsUpdateCache bool
	Op            types.Op
}

// Eval implements types.Expr
func (expr *Expr) Eval() float64 {
	if expr.ResultCache.Valid {
		return expr.ResultCache.Float64
	}
	var r float64
	if expr.Op == types.OpFetch {
		if expr.StaticValue.Valid {
			r = expr.StaticValue.Float64
		} else {
			r = expr.ValueLoader.Load()
		}
	} else {
		lhs := expr.LHS.Eval()
		rhs := expr.RHS.Eval()
		switch expr.Op {
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
		}
	}

	if expr.IsUpdateCache {
		expr.ResultCache.Float64 = r
		expr.ResultCache.Valid = true
	}

	return r
}

type stack []*Expr

func (s *stack) Push(expr Expr) *Expr {
	*s = append(*s, &expr)
	return s.First()
}

func (s *stack) First() *Expr {
	return (*s)[len(*s)-1]
}

func (s *stack) Pop() *Expr {
	r := s.First()
	*s = (*s)[:len(*s)-1]
	return r
}

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
//
// input example: "z x y + *"
// calculation interpretation: z * (x + y)
// tree: *(z,+(x,y))
func Parse(expression string, symResolver types.SymbolResolver) (*Expr, error) {

	parts := strings.Split(expression, " ")
	stack := stack{}
	for partIdx, part := range parts {
		if part == "" {
			continue
		}
		op := types.ParseOp(part)

		if op != types.OpUndefined {
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid expression '%s' at part index %d: expected at least two entries in the stack", expression, partIdx)
			}
			rhs := stack.Pop()
			lhs := stack.Pop()
			stack.Push(Expr{
				Symbol: part,
				LHS:    lhs,
				RHS:    rhs,
				Op:     op,
			})
			continue
		}

		valueLoader, err := types.ParseValue(part, symResolver)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value '%s': %w", part, err)
		}

		expr := stack.Push(Expr{
			Symbol:      part,
			ValueLoader: valueLoader,
			Op:          types.OpFetch,
		})
		if resolver, ok := valueLoader.(types.StaticValue); ok {
			expr.StaticValue.Valid = true
			expr.StaticValue.Float64 = resolver.Load()
		}
	}
	if len(stack) == 0 {
		return nil, fmt.Errorf("empty expression: '%s'", expression)
	}
	return stack[0], nil
}

// String implements types.Expr
func (expr Expr) String() string {
	switch expr.Op {
	case types.OpFetch:
		return expr.Symbol
	case types.OpIf:
		return fmt.Sprintf("(if %s>0 then %s)", expr.LHS, expr.RHS)
	default:
		return fmt.Sprintf("(%s %s %s)", expr.LHS, expr.Op.String(), expr.RHS)
	}
}

// ClearCache invalidates any cache is stored in the tree
func (expr *Expr) ClearCache() {
	if expr == nil {
		return
	}
	expr.ResultCache.Valid = false
	expr.LHS.ClearCache()
	expr.RHS.ClearCache()
}

// EnableUpdateCache defines if the cache should be set (when it is absent).
func (expr *Expr) EnableUpdateCache(newValue bool) {
	if expr == nil {
		return
	}
	expr.IsUpdateCache = newValue
	expr.LHS.EnableUpdateCache(newValue)
	expr.RHS.EnableUpdateCache(newValue)
}

// EnableMemoization implements types.Expr
func (expr *Expr) EnableMemoization(newValue bool) (oldValue bool) {
	oldValue = expr.IsUpdateCache
	expr.EnableUpdateCache(newValue)
	if !newValue {
		expr.ClearCache()
	}
	return
}
