package gometer

import "fmt"

// PrefixMetrics is a Metrics wrapper, that always add
// specified prefix to counters names.
type PrefixMetrics struct {
	Metrics
	prefix string
}

var _ Metrics = (*PrefixMetrics)(nil)

// Register calls underlying Metrics Register method with prefixed counterName.
func (m *PrefixMetrics) Register(counterName string, c *Counter) error {
	return m.Metrics.Register(m.prefix+counterName, c)
}

// Get calls underlying Metrics Get method with prefixed counterName.
func (m *PrefixMetrics) Get(counterName string) *Counter {
	return m.Metrics.Get(m.prefix + counterName)
}

// WithPrefix returns new PrefixMetrics with extended prefix.
func (m *PrefixMetrics) WithPrefix(prefix string, v ...interface{}) *PrefixMetrics {
	return &PrefixMetrics{
		Metrics: m.Metrics,
		prefix:  m.prefix + fmt.Sprintf(prefix, v...),
	}
}
