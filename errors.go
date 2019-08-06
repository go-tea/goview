package goview

// Error allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// IStatusError represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type IStatusError interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}
