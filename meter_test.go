package gometer

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsStartFileWriter(t *testing.T) {
	fileName := "test_write_to_file"
	metrics := New()
	metrics.SetFormatter(NewFormatter("\n"))

	require.Panics(t, func() {
		metrics.StartFileWriter(nil, FileWriterParams{})
	})

	inc := Counter{}
	inc.Add(10)
	err := metrics.Register("add_num", &inc)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	metrics.StartFileWriter(ctx, FileWriterParams{
		FilePath:       fileName,
		UpdateInterval: time.Second * 1,
	})

	testWriteToFile(t, testWriteToFileParams{
		fileName:      fileName,
		lineSeparator: "\n",
		expMetricCnt:  1,
		waitDur:       time.Second * 2,
		cancel:        cancel,
	})
	cancel()

	inc1 := Counter{}
	inc1.Add(4)
	err = metrics.Register("inc_num", &inc1)
	require.Nil(t, err)

	ctx, cancel = context.WithCancel(context.Background())
	metrics.StartFileWriter(ctx, FileWriterParams{
		FilePath:       fileName,
		UpdateInterval: time.Second * 2,
	})

	testWriteToFile(t, testWriteToFileParams{
		fileName:      fileName,
		lineSeparator: "\n",
		expMetricCnt:  2,
		waitDur:       time.Second * 3,
		cancel:        cancel,
	})
	cancel()
}

type testWriteToFileParams struct {
	fileName      string
	lineSeparator string

	expMetricCnt int

	waitDur time.Duration
	cancel  context.CancelFunc
}

func testWriteToFile(t *testing.T, p testWriteToFileParams) {
	time.Sleep(p.waitDur)

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

	c := Counter{}
	c.Add(10)
	err := metrics.Register("test_counter", &c)
	require.Nil(t, err)

	err = metrics.Write()
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

func TestMetricsDefault(t *testing.T) {
	SetOutput(newTestFile(t, "test_file"))

	SetFormatter(NewFormatter("\n"))
	assert.NotNil(t, std.formatter)

	c := Counter{}
	c.Add(10)
	err := Register("default_metrics_counter", &c)
	require.Nil(t, err)

	err = Write()
	require.Nil(t, err)
	err = os.Remove("test_file")
	require.Nil(t, err)

	SetErrorHandler(new(mockErrorHandler))
	assert.NotNil(t, std.errHandler)

	require.Panics(t, func() {
		StartFileWriter(nil, FileWriterParams{})
	})

	ctx, cancel := context.WithCancel(context.Background())
	StartFileWriter(ctx, FileWriterParams{
		FilePath:       "test_default_metrics",
		UpdateInterval: time.Minute,
	})
	cancel()

	counter := Get("default_metrics_counter")
	require.NotNil(t, counter)
	assert.Equal(t, int64(10), counter.Get())

	group := Group()
	assert.NotNil(t, group)

	err = RegisterGroup(group)
	assert.Nil(t, err)
}

type mockErrorHandler struct{}

func (e *mockErrorHandler) Handle(err error) {
	fmt.Fprintf(os.Stderr, "failed to write metrics file, %v\n", err)
}

func TestMetricsSetErrorHandler(t *testing.T) {
	metrics := New()
	metrics.SetErrorHandler(new(mockErrorHandler))
	assert.NotNil(t, metrics.errHandler)
}

func TestMetricsExistingCounter(t *testing.T) {
	metrics := New()
	counter := Counter{}
	err := metrics.Register("existing_metrics", &counter)
	require.Nil(t, err)

	err = metrics.Register("existing_metrics", &counter)
	require.NotNil(t, err)
}

func TestMetricsGetCounter(t *testing.T) {
	metrics := New()
	c := metrics.Get("not_existing_counter")
	require.Nil(t, c)

	counter := Counter{}
	counter.Set(10)
	err := metrics.Register("get_counter", &counter)
	require.Nil(t, err)

	c = metrics.Get("get_counter")
	require.NotNil(t, c)
	assert.Equal(t, int64(10), c.Get())
}

func TestMetricsGroup(t *testing.T) {
	metrics := New()

	group := metrics.Group()
	assert.NotNil(t, group)
}

func TestMetricsRegisterGroup(t *testing.T) {
	metrics := New()

	group := metrics.Group().WithPrefix("foo").WithSeparator(".")

	barCounter := Counter{}
	barCounter.Add(100)

	bazCounter := Counter{}
	bazCounter.Add(140)

	group.Add("bar", &barCounter)
	group.Add("baz", &bazCounter)

	err := metrics.RegisterGroup(group)
	require.Nil(t, err)

	gotBar := metrics.Get("foo.bar")
	require.NotNil(t, gotBar)
	assert.Equal(t, int64(100), gotBar.Get())
}
