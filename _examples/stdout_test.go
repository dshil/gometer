package example

import (
	"os"

	"github.com/dshil/gometer"
)

func ExampleWriteToStdout() {
	metric := gometer.New()

	metric.SetOutput(os.Stdout)
	metric.SetFormatter(gometer.NewDefaultFormatter())

	c := metric.NewCounter("http_requests_total")
	c.Add(1)

	metric.Write()
	// Output:
	// http_requests_total = 1
}
