package example

import (
	"fmt"
	"os"

	"github.com/dshil/gometer"
)

func ExampleWriteToStdout() {
	metrics := gometer.New()

	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(gometer.NewFormatter("\n"))

	c := metrics.NewCounter("http_requests_total")
	c.Add(1)

	if err := metrics.Write(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// Output:
	// http_requests_total = 1
}
