package gometer

import (
	"bytes"
	"fmt"
)

// SortedCounters represents counters slice sorted by name
type SortedCounters []struct {
	Name    string
	Counter *Counter
}

// Formatter determines a format of metrics representation.
type Formatter interface {
	Format(counters SortedCounters) []byte
}

// NewFormatter returns new default formatter.
//
// lineSeparator determines how one line of metric
// will be separated from another.
//
// As line separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default format for one line of metrics is: "%v = %v". Metrics will be sorted by key.
func NewFormatter(lineSeparator string) Formatter {
	return &defaultFormatter{
		lineSeparator: lineSeparator,
	}
}

type defaultFormatter struct {
	lineSeparator string
}

func (f *defaultFormatter) Format(counters SortedCounters) []byte {
	var buf bytes.Buffer

	for _, c := range counters {
		fmt.Fprintf(&buf, "%s = %d%s", c.Name, c.Counter.Get(), f.lineSeparator)
	}

	return buf.Bytes()
}

var _ Formatter = (*defaultFormatter)(nil)
