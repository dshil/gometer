package gometer

import "sync/atomic"

// Counter represents a kind of metric.
type Counter struct {
	val int64
}

// Add adds the corresponding value to a counter.
func (c *Counter) Add(val int64) {
	atomic.AddInt64(&c.val, val)
}

// Get returns the corresponding value for a counter.
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.val)
}

// Set sets the value to a counter. Value can be negative.
func (c *Counter) Set(val int64) {
	atomic.StoreInt64(&c.val, val)
}

// AddAndGet adds val to a counter and returns an updated value.
func (c *Counter) AddAndGet(val int64) int64 {
	return atomic.AddInt64(&c.val, val)
}
