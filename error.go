package gometer

// ErrorHandler is used to handle errors that can happen
// during async rewriting metrics file.
//
// Default error handler has the nil value.
// It will be a panic, if some errors will happen
// during async rewriting metrics file.
type ErrorHandler interface {
	// Handle handles the error for async rewriting metrics file.
	Handle(err error)
}
