package gometer

import "sync/atomic"

// Counter represents a general counter.
type Counter interface {
	// Add adds the corresponding value to counter.
	Add(val int64)

	// Set sets the value to counter. Value can be negative.
	Set(val int64)

	// Get returns the corresponding value for counter.
	Get() int64
}

type counter struct {
	val int64
}

func (c *counter) Add(val int64) {
	atomic.AddInt64(&c.val, val)
}

func (c *counter) Get() int64 {
	return atomic.LoadInt64(&c.val)
}

func (c *counter) Set(val int64) {
	atomic.StoreInt64(&c.val, val)
}
