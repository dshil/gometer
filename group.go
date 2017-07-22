package gometer

import "sync"

// CountersGroup represents a collection of grouped counters.
type CountersGroup struct {
	mu       sync.Mutex
	counters map[string]Counter
	prefix   string
}

// Add adds new counter in the group of counters. If a counter with `counterName` exists,
// it'll be overwritten.
func (g *CountersGroup) Add(counterName string, counter Counter) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counters[g.prefix+counterName] = counter
}

// Get returns a counter from the group of counters.
func (g *CountersGroup) Get(counterName string) Counter {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.counters[g.prefix+counterName]
}

// Counters returns a collection of grouped counters.
func (g *CountersGroup) Counters() map[string]Counter {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.counters
}
