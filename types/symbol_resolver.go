package types

// SymbolResolver is a dispatcher of variable names to their ValueLoader-s.
type SymbolResolver interface {
	// Resolve returns a ValueLoader for the variable of name `sym`.
	Resolve(sym string) (ValueLoader, error)
}
