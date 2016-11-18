package gometer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Metrics is a collection of counters.
type Metrics struct {
	mu        sync.Mutex
	out       io.Writer
	format    string
	counters  map[string]*Counter
	separator string
}

var std = New()

// New returns new emtpy collection of counters.
//
// out defines where to dump corresponding metrics.
// stderr is a default output destination.
//
// separator determines how one metric will be separated from another.
// default separator is a newline symbol.
//
// format determines how metric's name and metric's value
// will be dumped to output destination.
// default format is 'metric_name = metric_value'.
func New() *Metrics {
	m := &Metrics{
		out:       os.Stderr,
		counters:  make(map[string]*Counter),
		separator: "\n",
		format:    "%v = %v",
	}
	return m
}

// SetOutput sets output destination for standard metric.
// Default output destination is os.Stderr.
func (m *Metrics) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
}

// SetFormat sets printing format.
//
// Can be used any format that contains two values in a string representation.
// Each metric is a key-value pair.
//
// Examples of valid formats: "%v = %v"; "%v:%v"; "%v, %v".
// Examples of invalid formats: "%v"; "%v ="; "%v, ".
//
// Default format is: "%v = %v".
func (m *Metrics) SetFormat(f string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.format = f
}

// Write all existing metrics to output destination.
//
// Writing metrics to the file using this method will not recreate a file.
// it just append existing metrics to existing file's data.
// if you want to write metrics to clear file use WriteAtFile() method.
func (m *Metrics) Write() error {
	return write(m)
}

// Stopper allows to stop writing metrics to file.
type Stopper struct {
	Stop func()
}

// WriteToFile writes all metrics to clear file.
//
// updateInterval determines how often metric will be write to file.
// use stopper to stop writing metrics periodically to file.
func (m *Metrics) WriteToFile(path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	return writeToFile(m, path, updateInterval, runImmediately)
}

func writeToFile(m *Metrics, path string, updateInterval time.Duration, runImmediately bool) *Stopper {
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
	metric         *Metrics
	runImmediately bool
}

func runFileWriter(p fileWriterParams) {
	ticker := time.NewTicker(p.updateInterval)
	defer ticker.Stop()
	defer close(p.stopCh)

	if p.runImmediately {
		if err := createAndWriteFile(p.metric, p.path); err != nil {
			panic(err)
		}
	}

	for {
		select {
		case <-ticker.C:
			if err := createAndWriteFile(p.metric, p.path); err != nil {
				panic(err)
			}
		case <-p.stopCh:
			return
		}
	}
}

func createAndWriteFile(m *Metrics, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	m.SetOutput(file)
	return write(m)
}

func write(m *Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var buf bytes.Buffer
	for name, counter := range m.counters {
		metric := fmt.Sprintf(m.format, name, counter.Get()) + m.separator
		fmt.Fprint(&buf, metric)
	}

	if _, err := m.out.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// NewCounter creates new counter in metric collection and returns it.
func (m *Metrics) NewCounter(name string) *Counter {
	return newCounter(m, name)
}

func newCounter(m *Metrics, counterName string) *Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.counters[counterName]; ok {
		return c
	}

	c := &Counter{}
	m.counters[counterName] = c

	return c
}

// SetLineSeparator determines how one metric will be separated from another.
// As line separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default line separator is: "\n".
func (m *Metrics) SetLineSeparator(s string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.separator = s
}

// LineSeparator returns a line separator.
func (m *Metrics) LineSeparator() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.separator
}

// These functions are used for standard metric.

// SetOutput sets output destination for standard metric.
// Default output destination is os.Stderr.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
}

// SetFormat sets printing format for standard metric.
//
// Can be used any format that contains two values in string representation.
// Each metric is a key-value pair.
//
// Examples of valid formats: "%v = %v"; "%v:%v"; "%v, %v".
// Examples of invalid formats: "%v"; "%v ="; "%v, ".
//
// Default format is: "%v = %v".
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

// WriteToFile writes all metrics to clear file for standard metric.
//
// updateInterval determines how often metric will be write to file.
// use stopper to stop writing metrics periodically to file.
func WriteToFile(path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	return writeToFile(std, path, updateInterval, runImmediately)
}

// NewCounter returns new counter for standard metric.
func NewCounter(name string) *Counter {
	return newCounter(std, name)
}

// SetLineSeparator determines how one metric will be separated from another.
// As line separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default line separator is: "\n".
func SetLineSeparator(s string) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.separator = s
}

// LineSeparator returns metric's separator for standard metric.
func LineSeparator() string {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.separator
}
