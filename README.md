# gometer [![GoDoc](https://godoc.org/github.com/dshil/gometer?status.svg)](https://godoc.org/github.com/dshil/gometer) [![Build Status](https://travis-ci.org/dshil/gometer.svg?branch=master)](https://travis-ci.org/dshil/gometer) [![Go Report Card](https://goreportcard.com/badge/github.com/dshil/gometer)](https://goreportcard.com/report/github.com/dshil/gometer) [![Coverage Status](https://coveralls.io/repos/github/dshil/gometer/badge.svg)](https://coveralls.io/github/dshil/gometer) [![codebeat badge](https://codebeat.co/badges/e194755f-ceda-48dd-8c6d-dcfbca04e07b)](https://codebeat.co/projects/github-com-dshil-gometer)
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

	c := gometer.DefaultCounter{}
	c.Add(1)
	if err := metrics.Register("http_requests_total", &c); err != nil {
		fmt.Println(err.Error())
		return
	}

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

func (f *simpleFormatter) Format(counters map[string]gometer.Counter) []byte {
	var buf bytes.Buffer
	for name, counter := range counters {
		line := fmt.Sprintf("%v, %v", name, counter.Get()) + "\n"
		fmt.Fprint(&buf, line)
	}
	return buf.Bytes()
}

var _ gometer.Formatter = (*simpleFormatter)(nil)

func ExampleSimpleFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(new(simpleFormatter))

	c := gometer.DefaultCounter{}
	c.Add(100)
	if err := metrics.Register("http_requests_total", &c); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := metrics.Write(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// Output: http_requests_total, 100
}

type sortByNameFormatter struct{}

func (f *sortByNameFormatter) Format(counters map[string]gometer.Counter) []byte {
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

var _ gometer.Formatter = (*sortByNameFormatter)(nil)

func ExampleSortByNameFormatter() {
	metrics := gometer.New()
	metrics.SetOutput(os.Stdout)
	metrics.SetFormatter(new(sortByNameFormatter))

	adder := gometer.DefaultCounter{}
	adder.Add(10)
	if err := metrics.Register("adder", &adder); err != nil {
		fmt.Println(err.Error())
		return
	}

	setter := gometer.DefaultCounter{}
	setter.Set(-1)
	if err := metrics.Register("setter", &setter); err != nil {
		fmt.Println(err.Error())
		return
	}

	inc := gometer.DefaultCounter{}
	inc.Add(1)
	if err := metrics.Register("inc", &inc); err != nil {
		fmt.Println(err.Error())
		return
	}

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
	defer cancel()

	// write metrics to file periodically.
	gometer.StartFileWriter(ctx, gometer.FileWriterParams{
		FilePath:       "test_file",
		UpdateInterval: time.Second,
	})
}
```
