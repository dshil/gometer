package example

import (
	"os"

	"github.com/dshil/gometer"
)

func ExampleWriteToStdout() {
	metric := gometer.New()
	metric.SetOutput(os.Stdout)
	c := metric.NewCounter("num_counter")
	c.Add(1)

	metric.Write()
	// Output:
	// num_counter = 1
}
