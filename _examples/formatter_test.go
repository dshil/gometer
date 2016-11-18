package example

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/dshil/gometer"
)

type simpleFormatter struct{}

func (f *simpleFormatter) Format(counters map[string]*gometer.Counter) []byte {
	var buf bytes.Buffer
	for name, counter := range counters {
		line := fmt.Sprintf("%v, %v", name, counter.Get()) + "\n"
		fmt.Fprint(&buf, line)
	}
	return buf.Bytes()
}

func ExampleSimpleFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(new(simpleFormatter))

	c := metrics.NewCounter("http_requests_total")
	c.Add(100)

	if err := metrics.Write(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// Output: http_requests_total, 100
}

type sortByNameFormatter struct{}

func (f *sortByNameFormatter) Format(counters map[string]*gometer.Counter) []byte {
	var buf bytes.Buffer

	var names []string
	for name := range counters {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, n := range names {
		line := fmt.Sprintf("%v: %v", n, counters[n].Get()) + "\n"
		fmt.Fprintf(&buf, line)
	}

	return buf.Bytes()
}

func ExampleSortByNameFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(new(sortByNameFormatter))

	adder := metrics.NewCounter("adder")
	adder.Add(10)

	setter := metrics.NewCounter("setter")
	setter.Set(-1)

	inc := metrics.NewCounter("inc")
	inc.Add(1)

	if err := metrics.Write(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// Output:
	// adder: 10
	// inc: 1
	// setter: -1
}
