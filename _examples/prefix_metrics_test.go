package example

import (
	"fmt"
	"os"

	"github.com/dshil/gometer"
)

func ExamplePrefixMetrics() {
	prefixMetrics := gometer.New().WithPrefix("data.%s.%s.", "errors", "counters")
	prefixMetrics.SetOutput(os.Stdout)

	for _, name := range []string{"foo", "bar", "baz"} {
		c := prefixMetrics.Get(name)
		c.Add(100)
	}

	if err := prefixMetrics.Write(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// data.errors.counters.bar = 100
	// data.errors.counters.baz = 100
	// data.errors.counters.foo = 100
}
