package example

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dshil/gometer"
)

type simpleFormatter struct{}

func (f *simpleFormatter) Format(counters map[string]gometer.Counter) ([]byte, error) {
	var buf bytes.Buffer

	for name, counter := range counters {
		_, err := fmt.Fprintf(&buf, "%s:%d%s", name, counter.Get(), "\n")
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

var _ gometer.Formatter = (*simpleFormatter)(nil)

func ExampleSimpleFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(new(simpleFormatter))

	c := gometer.DefaultCounter{}
	c.Add(100)
	if err := metrics.Register("foo", &c); err != nil {
		fmt.Println(err)
		return
	}

	if err := metrics.Write(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// foo:100
}

func ExampleDefaultFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)

	for _, name := range []string{"foo", "bar", "baz"} {
		c := gometer.DefaultCounter{}
		c.Add(100)
		if err := metrics.Register(name, &c); err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := metrics.Write(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// bar = 100
	// baz = 100
	// foo = 100
}
