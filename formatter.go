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
	Format(counters map[string]Counter) []byte
}

// NewFormatter returns new default formatter.
//
// lineSeparator determines how one line of metric
// will be separated from another.
//
// As line separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default format for one line of metrics is: "%v = %v".
// defaultFormatter sorts metrics by value.
func NewFormatter(lineSeparator string) Formatter {
	df := &defaultFormatter{
		lineSeparator: lineSeparator,
	}
	return df
}

func sortedKeys(m map[string]Counter) []string {
	s := make([]string, 0, len(m))
	for key := range m {
		s = append(s, key)
	}
	sort.Strings(s)
	return s
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type defaultFormatter struct {
	lineSeparator string
}

func (f *defaultFormatter) Format(counters map[string]Counter) []byte {
	var buf bytes.Buffer

	for _, k := range sortedKeys(counters) {
		line := fmt.Sprintf("%v = %v", k, counters[k].Get()) + f.lineSeparator
		_, err := buf.WriteString(line)
		panicIfErr(err)
	}

	return buf.Bytes()
}

var _ Formatter = (*defaultFormatter)(nil)
