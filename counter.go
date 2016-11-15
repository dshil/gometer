package gometer

// Counter represents a general counter that allows to:
// increment, decrement, add, subtract, set the value.
type Counter interface {
	Incrementor

	// Set sets the value to counter. Value can be negative.
	Set(val int64)
}

type counter struct {
	value
}

func (c *counter) Add(val int64) {
	c.value.Add(val)
}

func (c *counter) Set(val int64) {
	c.value.Set(val)
}
