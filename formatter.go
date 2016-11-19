package gometer

import (
	"bytes"
	"fmt"
)

// Formatter is used to determine a format of metrics representation.
type Formatter interface {
	// Format is defined how metrics will be dumped
	// to output destination.
	Format(counters map[string]*Counter) []byte
}

type defaultFormatter struct {
	lineSeparator string
}

func (f *defaultFormatter) Format(counters map[string]*Counter) []byte {
	var buf bytes.Buffer
	for name, counter := range counters {
		line := fmt.Sprintf("%v = %v", name, counter.Get()) + f.lineSeparator
		fmt.Fprint(&buf, line)
	}
	return buf.Bytes()
}
