package example

import (
	"fmt"
	"time"

	"github.com/dshil/gometer"
)

func ExampleWriteToFile() {
	metrics := gometer.New()
	metrics.SetFormatter(gometer.NewFormatter("\n"))

	gometer.StartFileWriter(gometer.FileWriterParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	}).Stop()
}
