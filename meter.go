package gometer

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/dchest/safefile"
)

// Metrics is a collection of metrics.
type Metrics struct {
	wg       sync.WaitGroup
	stopOnce sync.Once
	cancelCh chan struct{}

	mu           sync.Mutex
	out          io.Writer
	counters     map[string]Counter
	formatter    Formatter
	panicHandler PanicHandler
}

// FileWriterParams represents a params for asynchronous file writing operation.
//
// FilePath represents a file path.
// UpdateInterval determines how often metrics data will be written to a file.
type FileWriterParams struct {
	FilePath       string
	UpdateInterval time.Duration
}

// Default is a standard metrics object.
var Default = New()

// New creates new empty collection of metrics.
//
// out defines where to write metrics.
// formatter determines how metric's values will be formatted.
func New() *Metrics {
	m := &Metrics{
		out:       os.Stderr,
		counters:  make(map[string]Counter),
		formatter: NewFormatter("\n"),
		cancelCh:  make(chan struct{}),
	}
	return m
}

// SetOutput sets output destination for metrics.
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

// Formatter returns a metrics formatter.
func (m *Metrics) Formatter() Formatter {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.formatter
}

// Register registers a new counter in metric collection, returns error if the counter
// with such name exists.
func (m *Metrics) Register(counterName string, c Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.counters[counterName]; ok {
		return fmt.Errorf("counter with name `%v` exists", counterName)
	}

	m.counters[counterName] = c
	return nil
}

// RegisterGroup registers a collection of counters in a metric collection, returns an
// error if a counter with such name exists.
func (m *Metrics) RegisterGroup(group *CountersGroup) error {
	for name, counter := range group.Counters() {
		if err := m.Register(name, counter); err != nil {
			return err
		}
	}
	return nil
}

// Get returns counter by name or nil if counter doesn't exist.
func (m *Metrics) Get(counterName string) Counter {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := m.counters[counterName]
	return c
}

// GetJSON filters counters by given predicate and returns them as a json
// marshaled map.
func (m *Metrics) GetJSON(predicate func(string) bool) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[string]Counter)
	formatter := jsonFormatter{}
	for k, v := range m.counters {
		if predicate(k) {
			result[k] = v
		}
	}

	data, err := formatter.Format(result)
	if err == nil {
		return data
	}

	if m.panicHandler != nil {
		m.panicHandler.Handle(err)
	} else {
		panic(err)
	}

	return nil
}

// SetPanicHandler sets error handler for errors that causing the panic.
func (m *Metrics) SetPanicHandler(handler PanicHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.panicHandler = handler
}

// Group returns a new counters group.
//
// During a registration a group of counters in the metrics collection, a base prefix
// will be added to each counter name in this group.
func (m *Metrics) Group(format string, v ...interface{}) *CountersGroup {
	return &CountersGroup{
		prefix:   fmt.Sprintf(format, v...),
		counters: make(map[string]Counter),
	}
}

// Write writes all existing metrics to output destination.
//
// Writing metrics to the file using this method will not recreate a file.
// It appends existing metrics to existing file's data.
// if you want to write metrics to clear file use WriteToFile() method.
func (m *Metrics) Write() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := m.formatter.Format(m.counters)
	if err != nil {
		return err
	}

	if _, err := m.out.Write(data); err != nil {
		return err
	}

	return nil
}

// StartFileWriter starts a goroutine that will periodically writes metrics to a file.
func (m *Metrics) StartFileWriter(p FileWriterParams) {
	m.wg.Add(1)

	go func() {
		defer m.wg.Done()

		ticker := time.NewTicker(p.UpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := m.createAndWriteFile(p.FilePath)
				if err != nil {
					if h := m.getPanicHandler(); h != nil {
						h.Handle(err)
						return
					}
					panic(err)
				}
			case <-m.cancelCh:
				return
			}
		}
	}()
}

// StopFileWriter stops a goroutine that will periodically writes metrics to a file.
func (m *Metrics) StopFileWriter() {
	m.stopOnce.Do(func() {
		close(m.cancelCh)
	})
	m.wg.Wait()
}

func (m *Metrics) getPanicHandler() PanicHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.panicHandler
}

// These functions are used for standard metrics.

// SetOutput sets output destination for standard metrics.
func SetOutput(out io.Writer) {
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.out = out
}

// SetFormatter sets formatter for standard metrics.
// Fore more details see Metrics.SetFormatter().
func SetFormatter(f Formatter) {
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.formatter = f
}

// Register registers a new counter in a metric collection, returns an error if a counter
// with such name exists.
func Register(counterName string, c Counter) error {
	return Default.Register(counterName, c)
}

// Get returns a counter by name or nil if the counter doesn't exist.
func Get(counterName string) Counter {
	return Default.Get(counterName)
}

// GetJSON filters counters by given predicate and returns them as a json
// marshaled map.
func GetJSON(predicate func(string) bool) []byte {
	return Default.GetJSON(predicate)
}

// SetPanicHandler sets error handler for errors that causing the panic.
func SetPanicHandler(handler PanicHandler) {
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.panicHandler = handler
}

// Write all existing metrics to an output destination.
// For more details see Metrics.Write().
func Write() error {
	return Default.Write()
}

// StartFileWriter writes all metrics to a clear file.
// For more details see Metrics.WriteToFile().
func StartFileWriter(p FileWriterParams) {
	Default.StartFileWriter(p)
}

// StopFileWriter stops a goroutine that will periodically writes metrics to a file.
func StopFileWriter() {
	Default.StopFileWriter()
}

// Group returns new group counter for the default metrics collection..
//
// During registration a group of counters in a metrics collection, a base prefix will be
// added to each counter name in this group.
func Group(format string, v ...interface{}) *CountersGroup {
	return &CountersGroup{
		prefix:   fmt.Sprintf(format, v...),
		counters: make(map[string]Counter),
	}
}

// RegisterGroup registers a collection of counters in a default metric collection,
// returns an error if a counter with such name exists.
func RegisterGroup(group *CountersGroup) error {
	return Default.RegisterGroup(group)
}

func (m *Metrics) createAndWriteFile(path string) error {
	// create an empty temporary file.
	file, err := safefile.Create(path, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	m.SetOutput(file)
	if err = m.Write(); err != nil {
		return err
	}

	// rename temporary file to existing.
	// it's necessary for atomic file rewriting.
	return file.Commit()
}
