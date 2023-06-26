package router

// Validator interface for types that need validation
type Validator interface {
	// Validate function validates based
	// on the state of the variable
	Validate() error
}
