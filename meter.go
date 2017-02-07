package gometer

import (
	"context"
	"fmt"
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

// Register registers new counter in metric collection, returns error if
// counter with such name exists.
func (m *Metrics) Register(counterName string, c *Counter) error {
	return registerCounter(m, counterName, c)
}

func registerCounter(metrics *Metrics, counterName string, counter *Counter) error {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	if _, ok := metrics.counters[counterName]; ok {
		return fmt.Errorf("counter with name `%v` exists", counterName)
	}

	metrics.counters[counterName] = counter
	return nil
}

// RegisterGroup registers a collection of counters in the metric collection, returns
// an error if a counter with such name exists.
func (m *Metrics) RegisterGroup(group *CountersGroup) error {
	return registerGroup(m, group)
}

func registerGroup(metrics *Metrics, group *CountersGroup) error {
	counters := group.Counters()

	for name, counter := range counters {
		if err := registerCounter(metrics, name, counter); err != nil {
			return err
		}
	}
	return nil
}

// Get returns counter by name or nil if counter doesn't exist.
func (m *Metrics) Get(counterName string) *Counter {
	return getCounter(m, counterName)
}

func getCounter(m *Metrics, counterName string) *Counter {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := m.counters[counterName]
	return c
}

// SetErrorHandler sets error handler for errors that
// can happen during async rewriting metrics file.
func (m *Metrics) SetErrorHandler(e ErrorHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errHandler = e
}

// Group returns new counters group.
//
// During the registration a group of counters in the metrics collection, the base
// prefix will be added to the each counter name in this group.
func (m *Metrics) Group(format string, v ...interface{}) *CountersGroup {
	return &CountersGroup{
		prefix:   fmt.Sprintf(format, v...),
		counters: make(map[string]*Counter),
	}
}

// Write all existing metrics to output destination.
//
// Writing metrics to the file using this method will not recreate a file.
// it appends existing metrics to existing file's data.
// if you want to write metrics to clear file use WriteToFile() method.
func (m *Metrics) Write() error {
	return write(m)
}

// FileWriterParams represents a params for async file writing operation.
//
// FilePath represents a file path.
// UpdateInterval determines how often metrics data will be dumped to file.
type FileWriterParams struct {
	FilePath       string
	UpdateInterval time.Duration
}

// StartFileWriter writes all metrics to clear file.
//
// updateInterval determines how often metric will be write to file.
func (m *Metrics) StartFileWriter(ctx context.Context, p FileWriterParams) {
	if ctx == nil {
		panic("nil Context")
	}
	go startFileWriter(ctx, m, p)
}

func startFileWriter(ctx context.Context, m *Metrics, p FileWriterParams) {
	ticker := time.NewTicker(p.UpdateInterval)
	defer ticker.Stop()

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

// Register registers new counter in metric collection, returns error if
// counter with such name exists.
func Register(counterName string, c *Counter) error {
	return registerCounter(std, counterName, c)
}

// Get returns counter by name or nil if counter doesn't exist.
func Get(counterName string) *Counter {
	return getCounter(std, counterName)
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

// StartFileWriter writes all metrics to clear file.
// For more details see Metrics.WriteToFile().
func StartFileWriter(ctx context.Context, p FileWriterParams) {
	if ctx == nil {
		panic("nil Context")
	}
	go startFileWriter(ctx, std, p)
}

// Group returns new group counter for the default metrics collection..
//
// During the registration a group of counters in the metrics collection, the base
// prefix will be added to each counter name in this group.
func Group(format string, v ...interface{}) *CountersGroup {
	return &CountersGroup{
		prefix:   fmt.Sprintf(format, v...),
		counters: make(map[string]*Counter),
	}
}

// RegisterGroup registers a collection of counters in the default metric collection,
// returns an error if a counter with such name exists.
func RegisterGroup(group *CountersGroup) error {
	return registerGroup(std, group)
}
