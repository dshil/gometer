package gometer

import (
	"bytes"
	"fmt"
	"sort"
)

// Formatter is used to determine a format of metrics representation.
type Formatter interface {
	// Format is defined how metrics will be dumped
	// to output destination.
	Format(counters map[string]*Counter) []byte
}

// NewFormatter returns new default formatter.
//
// lineSeparator determines how one line of metric
// will be separated from another.
//
// As line separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default format for one line of metrics is: "%v = %v".
// defaultFormatter sorts metrics by name.
func NewFormatter(lineSeparator string) Formatter {
	df := &defaultFormatter{
		lineSeparator: lineSeparator,
	}
	return df
}

type defaultFormatter struct {
	lineSeparator string
}

func (f *defaultFormatter) Format(counters map[string]*Counter) []byte {
	var buf bytes.Buffer

	var names []string
	for name := range counters {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, n := range names {
		line := fmt.Sprintf("%v = %v", n, counters[n].Get()) + f.lineSeparator
		fmt.Fprintf(&buf, line)
	}

	return buf.Bytes()
}
