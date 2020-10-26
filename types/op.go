package types

import (
	"fmt"
)

// Op is an identifier of a single operation of an expression.
type Op uint8

const (
	// See also a description of Reverse Polish Notation to better
	// understand the stack.

	// OpUndefined means operation was not successfully parsed
	OpUndefined = Op(iota)

	// OpFetch means put a value from the incoming values to the stack
	OpFetch

	// OpPlus means to add last two values from the stack, and put
	// the result to back the stack
	OpPlus

	// OpMinus means to subtract the last value from the stack from the
	// before last value, and put the result to back the stack.
	OpMinus

	// OpMinus means to multiply the last two value from the stack from the
	// before last value, and put the result to back the stack.
	OpMultiply

	// OpDivide means to divide the before-last value from the stack by
	// the last value from the stack, and put the result to back the stack.
	OpDivide

	// OpDivide means to take a power of the last value from the stack with
	// the base of the before-last value from the stack, and put the
	// result to back the stack.
	OpPower

	// OpIf means to take the last value from the stack if the before-last
	// value of the stack is greater than zero, and put it back to the stack.
	// If the before-last value is less or equals to zero, then to put
	// zero to the stack.
	OpIf

	// BoundaryOp could be used for iteration through all Op-s (to detect
	// the end of the iteration process).
	BoundaryOp
)

// String implements fmt.Stringer
func (op Op) String() string {
	switch op {
	case OpFetch:
		return "#"
	case OpPlus:
		return "+"
	case OpMinus:
		return "-"
	case OpMultiply:
		return "*"
	case OpDivide:
		return "/"
	case OpPower:
		return "^"
	case OpIf:
		return "if"
	default:
		return fmt.Sprintf("unknown_op_%d", op)
	}
}

// ParseOp returns an Op for a passed string Op name in `s`.
// It returns OpUndefined, if unable to parse.
func ParseOp(s string) Op {
	for opCandidate := OpFetch + 1; opCandidate < BoundaryOp; opCandidate++ {
		if s == opCandidate.String() {
			return opCandidate
		}
	}

	return OpUndefined
}
