package example

import (
	"os"

	"github.com/dshil/gometer"
)

func ExampleWriteToStdout() {
	metric := gometer.New()
	metric.SetOutput(os.Stdout)
	c := metric.NewCounter("num_counter")
	c.Inc()

	metric.Write()
	// Output:
	// num_counter = 1
}
