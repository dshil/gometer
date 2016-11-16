package gometer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type metric struct {
	mu             sync.Mutex
	out            io.Writer
	updateInterval time.Duration
	format         string
	metrics        map[string]Incrementor
	separator      string
}

var std = New()

// New returns new basic metric.
//
// out defines where to flush corresponding metric.
// stderr is a default output destination
//
// updateInterval defines how often metric will be flushed
// to output destination.
//
// updateInterval = 0 means that metric will not be flushed to
// output destination, in such case you need to flush metrics manually.
// you can use Write method for it.
//
// separator determines how one metric will be separated from another
// default separator is a newline symbol.
//
// format determines how format metric's name and metric's value
// default format is 'metric_name = metric_value'.
func New() *metric {
	m := &metric{
		out:       os.Stderr,
		metrics:   make(map[string]Incrementor),
		separator: "\n",
		format:    "%v = %v",
	}
	return m
}

// SetOutput sets output destination for metric.
func (m *metric) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
}

// SetUpdateInterval sets updateInterval for metric.
func (m *metric) SetUpdateInterval(t time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateInterval = t
}

// SetFormat sets printing format for metric.
func (m *metric) SetFormat(f string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.format = f
}

// Write all existing metrics to output destination for metric.
func (m *metric) Write() error {
	return write(m)
}

// WriteAtFile writes all metrics to clear file.
func (m *metric) WriteAtFile(path string) error {
	return writeAtFile(m, path)
}

// NewIncrementor returns new incrementor for metric.
func (m *metric) NewIncrementor(name string) Incrementor {
	return newIncrementor(m, name)
}

// NewCounter returns new counter for metric.
func (m *metric) NewCounter(name string) Counter {
	return newCounter(m, name)
}

// SetSeparator sets metric's separator for metric.
func (m *metric) SetSeparator(s string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.separator = s
}

// Separator returns metric's separator for metric.
func (m *metric) Separator() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.separator
}

func newIncrementor(m *metric, metricName string) Incrementor {
	m.mu.Lock()
	defer m.mu.Unlock()

	inc := &incrementor{
		value: value{},
	}

	m.metrics[metricName] = inc

	return inc
}

func newCounter(m *metric, metricName string) Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := &counter{
		value: value{},
	}

	m.metrics[metricName] = c

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

// Write all existing metrics to output destination for standard metric.
//
// Writing metrics to the file using this method will not recreate a file.
// it just append existing metrics to existing file's data.
// if you want to write metrics to clear file use WriteAtFile() method.
func Write() error {
	return write(std)
}

// WriteAtFile writes all metrics of standard metric to clear file.
func WriteAtFile(path string) error {
	return writeAtFile(std, path)
}

func writeAtFile(m *metric, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	m.SetOutput(file)
	return write(m)
}

func write(m *metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var buf bytes.Buffer
	for name, val := range m.metrics {
		fmt.Fprintf(&buf, m.format+m.separator, name, val.Value())
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

// SetSeparator sets metric's separator for standard metric.
func SetSeparator(s string) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.separator = s
}

// Separator returns metric's separator for standard metric.
func Separator() string {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.separator
}
