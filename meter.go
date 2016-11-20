package gometer

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/dchest/safefile"
)

// Metrics is a collection of metrics.
type Metrics struct {
	mu         sync.Mutex
	out        io.Writer
	counters   map[string]*Counter
	formatter  Formatter
	errHandler ErrorHandler
}

var std = New()

// New creates new empty collection of metrics.
//
// out defines where to dump corresponding metrics.
//
// formatterParams determines how metric's values
// will be divided one from another.
//
// As a formatter will be used default formatter with
// '\n' symbol as a metric line separator.
func New() *Metrics {
	m := &Metrics{
		out:       os.Stderr,
		counters:  make(map[string]*Counter),
		formatter: NewFormatter("\n"),
	}
	return m
}

// SetOutput sets output destination for metrics.
// Default output destination is os.Stderr.
func (m *Metrics) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
}

// SetFormatter sets a metrics's formatter.
func (m *Metrics) SetFormatter(f Formatter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.formatter = f
}

// Formatter returns metrics's formatter.
func (m *Metrics) Formatter() Formatter {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.formatter
}

// SetErrorHandler sets error handler for errors that
// can happen during async rewriting metrics file.
func (m *Metrics) SetErrorHandler(e ErrorHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errHandler = e
}

// Write all existing metrics to output destination.
//
// Writing metrics to the file using this method will not recreate a file.
// it appends existing metrics to existing file's data.
// if you want to write metrics to clear file use WriteToFile() method.
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
		metrics:        m,
		runImmediately: runImmediately,
	}
	go runFileWriter(params)
	return s
}

type fileWriterParams struct {
	stopCh         chan bool
	path           string
	updateInterval time.Duration
	metrics        *Metrics
	runImmediately bool
}

func runFileWriter(p fileWriterParams) {
	ticker := time.NewTicker(p.updateInterval)
	defer ticker.Stop()
	defer close(p.stopCh)

	if p.runImmediately {
		if err := createAndWriteFile(p.metrics, p.path); err != nil {
			if p.metrics.errHandler != nil {
				p.metrics.errHandler.Handle(err)
				return
			}
			panic(err)
		}
	}

	for {
		select {
		case <-ticker.C:
			if err := createAndWriteFile(p.metrics, p.path); err != nil {
				if p.metrics.errHandler != nil {
					p.metrics.errHandler.Handle(err)
					return
				}
				panic(err)
			}
		case <-p.stopCh:
			return
		}
	}
}

func createAndWriteFile(m *Metrics, path string) error {
	// create an empty temporary file.
	file, err := safefile.Create(path, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	m.SetOutput(file)
	if err = write(m); err != nil {
		return err
	}

	// rename temporary file to existing.
	// it's necessary for atomic file rewriting.
	if err = file.Commit(); err != nil {
		return err
	}

	return nil
}

func write(m *Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.out.Write(m.formatter.Format(m.counters)); err != nil {
		return err
	}
	return nil
}

// NewCounter creates new counter in metrics collection and returns it.
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

// These functions are used for standard metrics.

// SetOutput sets output destination for standard metrics.
// Default output destination is os.Stderr.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
}

// SetFormatter sets formatter for standard metrics.
// Fore more details see Metrics.SetFormatter().
func SetFormatter(f Formatter) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.formatter = f
}

// SetErrorHandler sets error handler for errors that
// can happen during async rewriting metrics file.
func SetErrorHandler(e ErrorHandler) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.errHandler = e
}

// Write all existing metrics to output destination.
// For more details see Metrics.Write().
func Write() error {
	return write(std)
}

// WriteToFile writes all metrics to clear file.
// For more details see Metrics.WriteToFile() .
func WriteToFile(path string, updateInterval time.Duration, runImmediately bool) *Stopper {
	return writeToFile(std, path, updateInterval, runImmediately)
}

// NewCounter returns new counter for standard metrics.
func NewCounter(name string) *Counter {
	return newCounter(std, name)
}
