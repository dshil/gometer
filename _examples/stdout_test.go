package example

import (
	"os"

	"github.com/dshil/gometer"
)

func ExampleWriteToStdout() {
	metric := gometer.New()
	metric.SetOutput(os.Stdout)
	c := metric.NewCounter(gometer.TotalHTTPRequests)
	c.Add(1)

	metric.Write()
	// Output:
	// http_requests_total = 1
}
