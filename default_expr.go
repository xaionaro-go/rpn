package rpn

import (
	tokenslice "github.com/xaionaro-go/rpn/implementations/tokenslice"
	"github.com/xaionaro-go/rpn/types"
)

// Expr is the default implementation of types.Expr.
//
// See also README.md.
type Expr = tokenslice.Expr

// Parse converts Reverse Polish Notation expression "expression" to
// a Eval()-uatable implementation Expr.
func Parse(expression string, symResolver types.SymbolResolver) (*Expr, error) {
	return tokenslice.Parse(expression, symResolver)
}
