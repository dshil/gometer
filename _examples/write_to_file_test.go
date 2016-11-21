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
	// call will stop writing to file operation.
	defer cancel()

	// write metrics to file periodically.
	gometer.WriteToFile(ctx, gometer.WriteToFileParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
		RunImmediately: true,
	})
}
