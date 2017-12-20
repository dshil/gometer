package example

import (
	"time"

	"github.com/dshil/gometer"
)

func ExampleWriteToFile() {
	metrics := gometer.New()
	metrics.SetFormatter(gometer.NewFormatter("\n"))

	gometer.StartFileWriter(gometer.FileWriterParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
	})
	gometer.StopFileWriter()
}
