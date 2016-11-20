package gometer

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsWriteToFile(t *testing.T) {
	fileName := "test_write_to_file"
	m := New()
	m.SetFormatter(NewFormatter("\n"))

	inc := m.NewCounter("add_num")
	inc.Add(10)

	stopper := m.WriteToFile(fileName, time.Second*10, true)
	testWriteToFile(t, testWriteToFileParams{
		fileName:      fileName,
		lineSeparator: "\n",
		expMetricCnt:  1,
		waitDur:       time.Second * 1,
		stopper:       stopper,
	})

	inc1 := m.NewCounter("inc_num")
	inc1.Add(4)

	stopper = m.WriteToFile(fileName, time.Second*10, true)
	testWriteToFile(t, testWriteToFileParams{
		fileName:      fileName,
		lineSeparator: "\n",
		expMetricCnt:  2,
		waitDur:       time.Second * 1,
		stopper:       stopper,
	})
}

type testWriteToFileParams struct {
	fileName      string
	lineSeparator string

	expMetricCnt int

	waitDur time.Duration
	stopper *Stopper
}

func testWriteToFile(t *testing.T, p testWriteToFileParams) {
	time.Sleep(p.waitDur)
	defer p.stopper.Stop()

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	err = os.Remove(p.fileName)
	require.Nil(t, err)

	metrics := strings.TrimSuffix(string(data), p.lineSeparator)
	metricsData := strings.Split(metrics, p.lineSeparator)
	require.Equal(t, p.expMetricCnt, len(metricsData))
}

func TestMetricsSetFormatter(t *testing.T) {
	fileName := "test_set_formatter"
	file := newTestFile(t, fileName)
	defer closeAndRemoveTestFile(t, file)

	metrics := New()
	metrics.SetOutput(file)
	metrics.SetFormatter(NewFormatter("\n"))

	c := metrics.NewCounter("test_counter")
	c.Add(10)

	err := metrics.Write()
	require.Nil(t, err)

	data, err := ioutil.ReadFile(fileName)
	require.Nil(t, err)
	metricsData := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	require.Equal(t, 1, len(metricsData))

	metricLine := strings.Split(metricsData[0], " = ")
	require.Equal(t, 2, len(metricLine))
	assert.Equal(t, "test_counter", metricLine[0])
	assert.Equal(t, "10", metricLine[1])
}

func TestMetricsFormatter(t *testing.T) {
	metrics := New()
	metrics.SetFormatter(NewFormatter("\n"))
	assert.NotNil(t, metrics.Formatter())
}

func newTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Create(fileName)
	require.Nil(t, err)
	return file
}

func closeAndRemoveTestFile(t *testing.T, f *os.File) {
	err := f.Close()
	require.Nil(t, err)
	err = os.Remove(f.Name())
	require.Nil(t, err)
}

func TestMetricsNewCounter(t *testing.T) {
	metrics := New()
	c1 := metrics.NewCounter("test_counter")

	// NewCounter will not recreate a counter (because of the same names),
	// just returns existing counter.
	c2 := metrics.NewCounter("test_counter")
	assert.Equal(t, c1, c2)
}

func TestMetricsDefault(t *testing.T) {
	SetOutput(os.Stderr)
	assert.Equal(t, os.Stderr, std.out)

	SetFormatter(NewFormatter("\n"))
	assert.NotNil(t, std.formatter)

	c := NewCounter("default_metrics_counter")
	require.NotNil(t, c)
	c.Add(10)

	err := Write()
	assert.Nil(t, err)

	SetErrorHandler(new(mockErrorHandler))
	assert.NotNil(t, std.errHandler)

	fileName := "default_metrics_file"
	stopper := WriteToFile(fileName, time.Second, true)
	require.NotNil(t, stopper)
	defer stopper.Stop()
}

type mockErrorHandler struct{}

func (e *mockErrorHandler) Handle(err error) {
	log.Printf("failed to write metrics file, %v\n", err)
}

func TestMetricsSetErrorHandler(t *testing.T) {
	metrics := New()
	metrics.SetErrorHandler(new(mockErrorHandler))
	assert.NotNil(t, metrics.errHandler)
}
