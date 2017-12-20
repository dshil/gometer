package gometer

// PanicHandler is used to handle errors that causing the panic.
type PanicHandler interface {
	Handle(err error)
}
