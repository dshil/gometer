package gometer

import "sync/atomic"

// Counter represents a one metric.
type Counter struct {
	val int64
}

// Add adds the corresponding value to counter.
func (c *Counter) Add(val int64) {
	atomic.AddInt64(&c.val, val)
}

// Get returns the corresponding value for counter.
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.val)
}

// Set sets the value to counter. Value can be negative.
func (c *Counter) Set(val int64) {
	atomic.StoreInt64(&c.val, val)
}
