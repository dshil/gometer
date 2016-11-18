package example

import (
	"fmt"
	"os"

	"github.com/dshil/gometer"
)

func ExampleSetFormatter() {
	formatFn := func(names ...interface{}) string {
		return fmt.Sprintf("%v = %v", names...)
	}

	m := gometer.New()
	c := m.NewCounter("test_counter")
	c.Set(1)
	m.SetFormatter(gometer.FormatterParams{
		LineSeparator: "\n",
		LineFormatter: formatFn,
	})
	m.SetOutput(os.Stdout)
	m.Write()
	// Output: test_counter = 1
}
