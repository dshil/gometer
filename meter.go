package gometer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Metric allows to add counters, incrementors,
// writes all metrics to out.
type Metric struct {
	mu             sync.Mutex
	out            io.Writer
	updateInterval time.Duration
	format         string
	metrics        map[string]Incrementor
}

var std = New(os.Stderr, 0)

// New returns new basic metric.
//
// out defines where to flush corresponding metric.
//
// updateInterval defines how often metric will be flushed to out.
//
// updateInterval = 0 means that metric will not be flushed to out,
// in such case you need to manually call Write method.
func New(out io.Writer, updateInterval time.Duration) *Metric {
	m := &Metric{
		out:            out,
		updateInterval: updateInterval,
		metrics:        make(map[string]Incrementor),
	}
	return m
}

// SetOutput sets output for metric.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
}

// SetUpdateInterval sets updateInterval for metric.
func SetUpdateInterval(t time.Duration) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.updateInterval = t
}

// SetFormat sets printing format to out.
func SetFormat(f string) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.format = f
}

// Write all existing metrics to io.Writer.
func Write() error {
	std.mu.Lock()
	defer std.mu.Unlock()

	var buf bytes.Buffer
	for name, val := range std.metrics {
		fmt.Fprintf(&buf, std.format, name, val.Value())
	}

	if _, err := std.out.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// NewIncrementor returns increment counter with 'name'.
func NewIncrementor(name string) Incrementor {
	std.mu.Lock()
	defer std.mu.Unlock()

	inc := &incrementor{
		value: value{},
	}

	std.metrics[name] = inc

	return inc
}

// NewCounter returns counter with 'name'.
func NewCounter(name string) Counter {
	std.mu.Lock()
	defer std.mu.Unlock()

	c := &counter{
		value: value{},
	}

	std.metrics[name] = c

	return c
}
