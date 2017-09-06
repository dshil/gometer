package example

import (
	"time"

	"github.com/dshil/gometer"
)

func ExampleWriteToFile() {
	metrics := gometer.New()
	metrics.SetFormatter(gometer.NewFormatter("\n"))

	// write metrics to file periodically.
	gometer.StartFileWriter(gometer.FileWriterParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
	})
	gometer.StopFileWriter()
}
