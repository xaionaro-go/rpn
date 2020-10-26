package types

// NullFloat64 is just a float64 with a bool which could be used
// to mark if the value is set.
type NullFloat64 struct {
	// Float64 is the value.
	Float64 float64

	// Valid is the marker if the value is set.
	Valid bool
}
