package types

import (
	"fmt"
)

// Expr is a parsed expression which could be executed by method Eval
// and will return the calculated value.
type Expr interface {
	// Eval executes the expression
	Eval() float64

	// EnableMemoization defines if memoization (caching of resulting
	// values) should be used
	EnableMemoization(bool) bool

	fmt.Stringer
}
