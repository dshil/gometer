package gometer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type metric struct {
	mu           sync.Mutex
	out          io.Writer
	format       string
	counters     map[string]Counter
	incrementors map[string]Incrementor
	separator    string
}

var std = New()

// New returns new basic metric.
//
// out defines where to flush corresponding metric.
// stderr is a default output destination
//
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
		out:          os.Stderr,
		counters:     make(map[string]Counter),
		incrementors: make(map[string]Incrementor),
		separator:    "\n",
		format:       "%v = %v",
	}
	return m
}

// SetOutput sets output destination for metric.
func (m *metric) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
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

// Stopper allows to stop writing metrics to file.
type Stopper struct {
	Stop func()
}

// WriteAtFile writes all metrics to clear file.
//
// updateInterval determines how often metric will be write
// to file.
// use Stopper to stop writing metrics periodically to file.
func (m *metric) WriteAtFile(path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	return writeAtFile(m, path, updateInterval, runImmediately)
}

func writeAtFile(m *metric, path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	stopCh := make(chan bool, 1)

	once := sync.Once{}
	s := &Stopper{
		Stop: func() {
			once.Do(func() {
				stopCh <- true
			})
		},
	}

	params := fileWriterParams{
		stopCh:         stopCh,
		path:           path,
		updateInterval: updateInterval,
		metric:         m,
		runImmediately: runImmediately,
	}
	go runFileWriter(params)
	return s
}

type fileWriterParams struct {
	stopCh         chan bool
	path           string
	updateInterval time.Duration
	metric         *metric
	runImmediately bool
}

func runFileWriter(p fileWriterParams) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("faile to write a file %v, recovered, error: %v\n", p.path, e)
		}
	}()

	if p.runImmediately {
		if err := createAndWriteFile(p.metric, p.path); err != nil {
			panic(err)
		}
	}

	ticker := time.NewTicker(p.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := createAndWriteFile(p.metric, p.path); err != nil {
				panic(err)
			}
		case <-p.stopCh:
			close(p.stopCh)
			return
		}
	}
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

	m.incrementors[metricName] = inc

	return inc
}

func newCounter(m *metric, metricName string) Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := &counter{
		value: value{},
	}

	m.counters[metricName] = c

	return c
}

// SetOutput sets output destination for standard metric.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
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
// it writes to file periodically, until you don't stop it.
func WriteAtFile(path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	return writeAtFile(std, path, updateInterval, runImmediately)
}

func createAndWriteFile(m *metric, path string) error {
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
	for name, val := range m.counters {
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
