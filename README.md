# gometer [![GoDoc](https://godoc.org/github.com/dshil/gometer?status.svg)](https://godoc.org/github.com/dshil/gometer) [![Build Status](https://travis-ci.org/dshil/gometer.svg?branch=master)](https://travis-ci.org/dshil/gometer) [![Coverage Status](https://coveralls.io/repos/github/dshil/gometer/badge.svg)](https://coveralls.io/github/dshil/gometer)


`gometer` is a small library for your application's metrics.

Basically, the main goal of `gometer` is to represent key-value metrics in some format.   
Later these formatted metrics can be used by other services (e.g. Zabbix).

## Installation

Install [Go](https://golang.org/) and run:

    go get -v github.com/dshil/gometer


## Documentation

Documentation is available on [GoDoc](https://godoc.org/github.com/dshil/gometer).

## Quick start

##### Write metrics to stdout.

```go
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
```

##### Own formatter for metrics representation.

```go
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
```

##### Write metrics to file periodically with cancelation.

```go
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
```
