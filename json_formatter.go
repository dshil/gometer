package gometer

import (
	"bytes"
	"fmt"
)

type jsonFormatter struct {
}

func (f *jsonFormatter) Format(counters map[string]Counter) ([]byte, error) {
	var buf bytes.Buffer

	_, err := buf.WriteRune('{')
	if err != nil {
		return nil, err
	}

	first := true
	for _, k := range sortedKeys(counters) {
		if first {
			first = false
		} else {
			if _, err := buf.WriteRune(','); err != nil {
				return nil, err
			}
		}
		if _, err := fmt.Fprintf(&buf, `"%s":%d`, k, counters[k].Get()); err != nil {
			return nil, err
		}
	}

	if _, err := buf.WriteRune('}'); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var _ Formatter = (*jsonFormatter)(nil)
