package gometer

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/dchest/safefile"
)

// Metrics is a collection of metrics.
type Metrics interface {
	SetOutput(io.Writer)
	SetFormatter(Formatter)
	Formatter() Formatter
	Get(string) *Counter
	GetJSON(func(string) bool) []byte
	WithPrefix(string, ...interface{}) *PrefixMetrics
	Write() error
	StartFileWriter(FileWriterParams) Stopper
}

// DefaultMetrics is a default implementation of Metrics.
type DefaultMetrics struct {
	wg       sync.WaitGroup
	stopOnce sync.Once
	cancelCh chan struct{}

	mu         sync.Mutex
	out        io.Writer
	counters   map[string]*Counter
	formatter  Formatter
	rootPrefix string
}

var _ Metrics = (*DefaultMetrics)(nil)
var _ Metrics = (*PrefixMetrics)(nil)

// FileWriterParams represents a params for asynchronous file writing operation.
//
// FilePath represents a file path.
// UpdateInterval determines how often metrics data will be written to a file.
// NoFlushOnStop disables metrics flushing when the metrics writer finishes.
// ErrorHandler allows to handle errors from the goroutine that writes metrics.
type FileWriterParams struct {
	FilePath       string
	UpdateInterval time.Duration
	NoFlushOnStop  bool
	ErrorHandler   func(err error)
}

// Default is a standard metrics object.
var Default = New()

// New creates new empty collection of metrics.
func New() *DefaultMetrics {
	m := &DefaultMetrics{
		out:       os.Stderr,
		counters:  make(map[string]*Counter),
		formatter: NewFormatter("\n"),
		cancelCh:  make(chan struct{}),
	}
	return m
}

// SetOutput sets output destination for metrics.
func (m *DefaultMetrics) SetOutput(out io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.out = out
}

// SetFormatter sets a metrics's formatter.
func (m *DefaultMetrics) SetFormatter(f Formatter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.formatter = f
}

// SetRootPrefix sets root prefix used to format output
func (m *DefaultMetrics) SetRootPrefix(prefix string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rootPrefix = prefix
}

// Formatter returns a metrics formatter.
func (m *DefaultMetrics) Formatter() Formatter {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.formatter
}

// Get returns counter by name. If counter doesn't exist it will be created.
func (m *DefaultMetrics) Get(counterName string) *Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.counters[counterName]; ok {
		return c
	}

	c := &Counter{}
	m.counters[counterName] = c
	return c
}

// GetJSON filters counters by given predicate and returns them as a json marshaled map.
func (m *DefaultMetrics) GetJSON(predicate func(string) bool) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[string]*Counter)
	for k, v := range m.counters {
		if predicate(k) {
			result[k] = v
		}
	}

	formatter := jsonFormatter{}
	return formatter.Format(m.makeSortedCounters(result))
}

func (m *DefaultMetrics) makeSortedCounters(counters map[string]*Counter) SortedCounters {
	s := make(SortedCounters, 0, len(counters))
	for k, v := range counters {
		s = append(s, struct {
			Name    string
			Counter *Counter
		}{Name: m.rootPrefix + k, Counter: v})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name < s[j].Name
	})
	return s
}

// Write writes all existing metrics to output destination.
//
// Writing metrics to the file using this method will not recreate a file.
// It appends existing metrics to existing file's data.
// if you want to write metrics to clear file use StartFileWriter() method.
func (m *DefaultMetrics) Write() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := m.formatter.Format(m.makeSortedCounters(m.counters))

	if _, err := m.out.Write(data); err != nil {
		return err
	}

	return nil
}

// StartFileWriter starts a goroutine that periodically writes metrics to a file.
func (m *DefaultMetrics) StartFileWriter(params FileWriterParams) Stopper {
	m.wg.Add(1)

	go func() {
		defer m.wg.Done()
		m.run(params)
	}()

	return &stopperFunc{stop: func() {
		m.stopOnce.Do(func() {
			close(m.cancelCh)

			m.wg.Wait()
			if !params.NoFlushOnStop {
				m.handleFileWrite(params)
			}
		})
	}}
}

// WithPrefix creates new PrefixMetrics that uses original Metrics with specified prefix.
func (m *DefaultMetrics) WithPrefix(prefix string, v ...interface{}) *PrefixMetrics {
	return &PrefixMetrics{
		Metrics: m,
		prefix:  fmt.Sprintf(prefix, v...),
	}
}

// These functions are used for standard metrics.

// SetOutput sets output destination for standard metrics.
func SetOutput(out io.Writer) {
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.out = out
}

// SetFormatter sets formatter for standard metrics.
// Fore more details see DefaultMetrics.SetFormatter().
func SetFormatter(f Formatter) {
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.formatter = f
}

// Get returns counter by name. If counter doesn't exist it will be created.
func Get(counterName string) *Counter {
	return Default.Get(counterName)
}

// GetJSON filters counters by given predicate and returns them as a json marshaled map.
func GetJSON(predicate func(string) bool) []byte {
	return Default.GetJSON(predicate)
}

// Write all existing metrics to an output destination.
// For more details see DefaultMetrics.Write().
func Write() error {
	return Default.Write()
}

// StartFileWriter starts a goroutine that periodically writes metrics to a file.
// For more details see DefaultMetrics.StartFileWriter().
func StartFileWriter(p FileWriterParams) Stopper {
	return Default.StartFileWriter(p)
}

// WithPrefix creates new PrefixMetrics that uses original Metrics with specified prefix.
// For more details see DefaultMetrics.WithPrefix().
func WithPrefix(prefix string, v ...interface{}) *PrefixMetrics {
	return Default.WithPrefix(prefix, v...)
}

func (m *DefaultMetrics) run(params FileWriterParams) {
	ticker := time.NewTicker(params.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.handleFileWrite(params)
		case <-m.cancelCh:
			return
		}
	}
}

func (m *DefaultMetrics) handleFileWrite(params FileWriterParams) {
	err := m.createAndWriteFile(params.FilePath)
	if err != nil {
		if params.ErrorHandler != nil {
			params.ErrorHandler(err)
		} else {
			panic(err)
		}
	}
}

func (m *DefaultMetrics) createAndWriteFile(path string) error {
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
