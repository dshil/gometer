package gometer

import (
	"bytes"
	"fmt"
)

type jsonFormatter struct {
}

func (f *jsonFormatter) Format(counters map[string]Counter) []byte {
	var buf bytes.Buffer

	_, err := buf.WriteRune('{')
	panicIfErr(err)
	first := true
	for _, k := range sortedKeys(counters) {
		if first {
			first = false
		} else {
			_, err := buf.WriteRune(',')
			panicIfErr(err)
		}
		_, err := fmt.Fprintf(&buf, `"%s":%v`, k, counters[k].Get())
		panicIfErr(err)
	}
	_, err = buf.WriteRune('}')
	panicIfErr(err)
	return buf.Bytes()
}

var _ Formatter = (*jsonFormatter)(nil)
