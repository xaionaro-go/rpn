package rpn

import (
	"strings"

	callslice "github.com/xaionaro-go/rpn/implementations/callslice"
	calltree "github.com/xaionaro-go/rpn/implementations/calltree"
	"github.com/xaionaro-go/rpn/types"
)

// Expr is just an interface for a parsed expression.
//
// See also README.md.
type Expr = types.Expr

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
func Parse(expression string, symResolver types.SymbolResolver) (Expr, error) {
	if len(strings.Split(expression, " ")) > 20 {
		return callslice.Parse(expression, symResolver)
	}
	return calltree.Parse(expression, symResolver)
}
