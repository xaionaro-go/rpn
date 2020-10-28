package types

// StaticValue is an implementation of ValueLoader which is just
// a static float64 value. Using of this type allows to avoid extra
// function calls, sometimes.
type StaticValue float64

// Load implements ValueLoader.
//go:nosplit
func (r StaticValue) Load() float64 {
	return (float64)(r)
}

// FuncValue is just a function wrapper which implements ValueLoader.
type FuncValue func() float64

// Load implements ValueLoader.
//go:nosplit
func (r FuncValue) Load() float64 {
	return r()
}

// ValueLoader is something able to return a value of the variable.
type ValueLoader interface {
	// Load returns the value of the variable.
	Load() float64
}
