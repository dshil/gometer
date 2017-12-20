package gometer

import (
	"bytes"
	"fmt"
	"sort"
)

// Formatter determines a format of metrics representation.
type Formatter interface {
	Format(counters map[string]Counter) ([]byte, error)
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

func (f *defaultFormatter) Format(counters map[string]Counter) ([]byte, error) {
	var buf bytes.Buffer

	for _, k := range sortedKeys(counters) {
		_, err := fmt.Fprintf(&buf, "%s = %d%s", k, counters[k].Get(), f.lineSeparator)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

var _ Formatter = (*defaultFormatter)(nil)

func sortedKeys(m map[string]Counter) []string {
	s := make([]string, 0, len(m))
	for key := range m {
		s = append(s, key)
	}
	sort.Strings(s)
	return s
}
