package gometer

import (
	"fmt"
	"sync"
)

// GroupCounter represents a collection of grouped counters.
type GroupCounter struct {
	mu              sync.Mutex
	counters        map[string]*Counter
	prefix          string
	prefixSeparator string
}

// WithPrefix returns a group counter with a base prefix.
//
// During the registration a group of counters in the metrics collection, the base
// prefix will be added to each counter name in this group.
func (g *GroupCounter) WithPrefix(format string, v ...interface{}) *GroupCounter {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prefix = fmt.Sprintf(format, v...)

	return g
}

// WithSeparator returns a group counter with a prefix separator.
//
// prefixSeparator determines how base prefix will be separated from the counter name.
func (g *GroupCounter) WithSeparator(prefixSeparator string) *GroupCounter {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prefixSeparator = prefixSeparator

	return g
}

// Separator returns a prefix separator for the group.
func (g *GroupCounter) Separator() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.prefixSeparator
}

// Prefix returns a base prefix for a group counter.
func (g *GroupCounter) Prefix() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.prefix
}

// Add adds new counter in the group of counters.
func (g *GroupCounter) Add(counterName string, counter *Counter) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counters[counterName] = counter
}

// Get returns a counter from the group of counters.
func (g *GroupCounter) Get(counterName string) *Counter {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.counters[counterName]
}

// Counters returns a collection of grouped counters with a formatted names.
//
// The name of the formatted counter will be constructed as a group prefix + a prefix
// separator + an actual counter name.
func (g *GroupCounter) Counters() map[string]*Counter {
	g.mu.Lock()
	defer g.mu.Unlock()

	counters := make(map[string]*Counter, len(g.counters))
	for name, counter := range g.counters {
		counters[g.prefix+g.prefixSeparator+name] = counter
	}

	return counters
}
