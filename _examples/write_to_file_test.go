package example

import (
	"context"
	"time"

	"github.com/dshil/gometer"
)

func ExampleWriteToFile() {
	metrics := gometer.New()
	metrics.SetFormatter(gometer.NewFormatter("\n"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// write metrics to file periodically.
	gometer.StartFileWriter(ctx, gometer.FileWriterParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
	})
}
