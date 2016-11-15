package gometer

import "errors"

// Incrementor represents an increment counter.
type Incrementor interface {
	// Inc adds 1 to counter.
	Inc()

	// Add adds the corresponding value to counter.
	// It will panic if val < 0.
	Add(val int64)

	// Value returns current value of counter.
	Value() int64
}

type incrementor struct {
	value
}

func (i *incrementor) Inc() {
	i.value.Inc()
}

func (i *incrementor) Add(val int64) {
	if val < 0 {
		panic(errors.New("counter can only increase own value"))
	}
	i.value.Add(val)
}
