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

// SetOutput sets output destination for metric.
func (m *Metric) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
}

// SetUpdateInterval sets updateInterval for metric.
func (m *Metric) SetUpdateInterval(t time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateInterval = t
}

// SetFormat sets printing format for metric.
func (m *Metric) SetFormat(f string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.format = f
}

func (m *Metric) Write() error {
	return write(m)
}

// NewIncrementor returns new incrementor for metric.
func (m *Metric) NewIncrementor(name string) Incrementor {
	return newIncrementor(m, name)
}

// NewCounter returns new counter for metric.
func (m *Metric) NewCounter(name string) Counter {
	return newCounter(m, name)
}

func newIncrementor(m *Metric, metricName string) Incrementor {
	m.mu.Lock()
	defer m.mu.Unlock()

	inc := &incrementor{
		value: value{},
	}

	m.metrics[metricName] = inc

	return inc
}

func newCounter(m *Metric, metricName string) Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := &counter{
		value: value{},
	}

	std.metrics[metricName] = c

	return c
}

// SetOutput sets output destination for standard metric.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
}

// SetUpdateInterval sets updateInterval for standard metric.
func SetUpdateInterval(t time.Duration) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.updateInterval = t
}

// SetFormat sets printing format for standard metric.
func SetFormat(f string) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.format = f
}

// Write all existing metrics output destination for standard metric.
func Write() error {
	return write(std)
}

func write(m *Metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var buf bytes.Buffer
	for name, val := range std.metrics {
		fmt.Fprintf(&buf, m.format, name, val.Value())
	}

	if _, err := m.out.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// NewIncrementor returns new incrementor for standard metric.
func NewIncrementor(name string) Incrementor {
	return newIncrementor(std, name)
}

// NewCounter returns new counter for standard metric.
func NewCounter(name string) Counter {
	return newCounter(std, name)
}
