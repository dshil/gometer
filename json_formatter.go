package gometer

import (
	"bytes"
	"fmt"
)

type jsonFormatter struct {
}

// Format implements Formatter.Format.
// It always returns nil error.
func (f *jsonFormatter) Format(counters map[string]Counter) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteRune('{')

	first := true
	for _, k := range sortedKeys(counters) {
		if first {
			first = false
		} else {
			buf.WriteRune(',')
		}
		fmt.Fprintf(&buf, `"%s":%d`, k, counters[k].Get())
	}

	buf.WriteRune('}')

	return buf.Bytes(), nil
}

var _ Formatter = (*jsonFormatter)(nil)
