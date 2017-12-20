# gometer [![GoDoc](https://godoc.org/github.com/dshil/gometer?status.svg)](https://godoc.org/github.com/dshil/gometer) [![Build Status](https://travis-ci.org/dshil/gometer.svg?branch=master)](https://travis-ci.org/dshil/gometer) [![Go Report Card](https://goreportcard.com/badge/github.com/dshil/gometer)](https://goreportcard.com/report/github.com/dshil/gometer) [![Coverage Status](https://coveralls.io/repos/github/dshil/gometer/badge.svg)](https://coveralls.io/github/dshil/gometer)
`gometer` is a small library for your application's metrics.

The main goal of `gometer` is to be very simple, small and stupid, and to write
formatted key-value metrics somewhere.

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

##### Write metrics to a file periodically.

```go
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
```

##### Own formatter for metrics representation.

```go
package example

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dshil/gometer"
)

type simpleFormatter struct{}

func (f *simpleFormatter) Format(counters map[string]gometer.Counter) []byte {
	var buf bytes.Buffer

	for name, counter := range counters {
		fmt.Fprintf(&buf, "%s:%d%s", name, counter.Get(), "\n")
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
```

##### Get metrics in JSON format for a specified pattern.

```go
package example

import (
	"fmt"

	"github.com/dshil/gometer"
	"github.com/gobwas/glob"
)

func ExampleMetricsGetJSONGlobPatterns() {
	metrics := gometer.New()

	for k, v := range map[string]int64{
		"abc":  10,
		"abb":  42,
		"adc":  33,
		"aaac": 17,
	} {
		c := new(gometer.DefaultCounter)
		c.Set(v)
		if err := metrics.Register(k, c); err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, tCase := range [...]struct {
		pattern  string
		expected string
	}{
		{
			pattern:  "*",
			expected: `{"abc": 10, "abb": 42, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "a*",
			expected: `{"abc": 10, "abb": 42, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "a?c",
			expected: `{"abc": 10, "adc": 33}`,
		},
		{
			pattern:  "a*c",
			expected: `{"abc": 10, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "*b*",
			expected: `{"abc": 10, "abb": 42}`,
		},
		{
			pattern:  "??[ab]*",
			expected: `{"abb": 42, "aaac":17}`,
		},
	} {
		g := glob.MustCompile(tCase.pattern)
		b := metrics.GetJSON(g.Match)
		fmt.Println(string(b))
	}

	// Output:
	// {"aaac":17,"abb":42,"abc":10,"adc":33}
	// {"aaac":17,"abb":42,"abc":10,"adc":33}
	// {"abc":10,"adc":33}
	// {"aaac":17,"abc":10,"adc":33}
	// {"abb":42,"abc":10}
	// {"aaac":17,"abb":42}
}
```
