package gometer

import (
	"context"
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

// WriteToFileParams represents a params for async file writing operation.
//
// FilePath represents a file path.
// UpdateInterval determines how often metrics data will be dumped to file.
// RunImmediately allows to immediatelly dump all metrics data to file.
type WriteToFileParams struct {
	FilePath       string
	UpdateInterval time.Duration
	RunImmediately bool
}

// WriteToFile writes all metrics to clear file.
//
// updateInterval determines how often metric will be write to file.
func (m *Metrics) WriteToFile(ctx context.Context, p WriteToFileParams) {
	go runFileWriter(ctx, m, p)
}

func runFileWriter(ctx context.Context, m *Metrics, p WriteToFileParams) {
	ticker := time.NewTicker(p.UpdateInterval)
	defer ticker.Stop()

	if p.RunImmediately {
		if err := createAndWriteFile(m, p.FilePath); err != nil {
			if m.errHandler != nil {
				m.errHandler.Handle(err)
				return
			}
			panic(err)
		}
	}

	for {
		select {
		case <-ticker.C:
			if err := createAndWriteFile(m, p.FilePath); err != nil {
				if m.errHandler != nil {
					m.errHandler.Handle(err)
					return
				}
				panic(err)
			}
		case <-ctx.Done():
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
func WriteToFile(ctx context.Context, p WriteToFileParams) {
	go runFileWriter(ctx, std, p)
}

// NewCounter returns new counter for standard metrics.
func NewCounter(name string) *Counter {
	return newCounter(std, name)
}
