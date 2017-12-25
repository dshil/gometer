package gometer

import (
	"bytes"
	"fmt"
)

type jsonFormatter struct {
}

func (f *jsonFormatter) Format(counters SortedCounters) []byte {
	var buf bytes.Buffer

	buf.WriteRune('{')

	first := true
	for _, c := range counters {
		if first {
			first = false
		} else {
			buf.WriteRune(',')
		}
		fmt.Fprintf(&buf, `"%s":%d`, c.Name, c.Counter.Get())
	}

	buf.WriteRune('}')

	return buf.Bytes()
}

var _ Formatter = (*jsonFormatter)(nil)
