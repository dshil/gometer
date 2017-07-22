package gometer

import "sync/atomic"

// Counter represents a kind of metric.
type Counter interface {
	// Add adds the corresponding value to a counter.
	Add(v int64)

	// Get returns the corresponding value for a counter.
	Get() int64

	// Set sets the value to a counter. Value can be negative.
	Set(v int64)
}

// DefaultCounter implements Counter.
type DefaultCounter struct {
	val int64
}

// Add implements Counter.Add
func (c *DefaultCounter) Add(val int64) {
	atomic.AddInt64(&c.val, val)
}

// Get implements Counter.Get
func (c *DefaultCounter) Get() int64 {
	return atomic.LoadInt64(&c.val)
}

// Set implements Counter.Set
func (c *DefaultCounter) Set(val int64) {
	atomic.StoreInt64(&c.val, val)
}

// AddAndGet adds val to a counter and returns an updated value.
func (c *DefaultCounter) AddAndGet(val int64) int64 {
	return atomic.AddInt64(&c.val, val)
}
