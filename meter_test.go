package gometer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsStopWithoutStart(t *testing.T) {
	t.Parallel()

	metrics := New()
	metrics.StopFileWriter()
}

func TestMetricsStopTwice(t *testing.T) {
	t.Parallel()

	metrics := New()
	metrics.StartFileWriter(FileWriterParams{
		FilePath:       os.DevNull,
		UpdateInterval: time.Millisecond * 100,
	})
	metrics.StopFileWriter()
	metrics.StopFileWriter()
}

func TestMetricsStartFileWriter(t *testing.T) {
	t.Parallel()

	file := newTempFile(t)
	require.Nil(t, file.Close())
	defer os.Remove(file.Name())

	metrics := New()
	lineSep := "\n"

	inc := metrics.Get("add_num")
	inc.Add(10)

	metrics.StartFileWriter(FileWriterParams{
		FilePath:       file.Name(),
		UpdateInterval: time.Millisecond * 100,
	})
	defer metrics.StopFileWriter()

	checkFileWriter(t, file.Name(), lineSep, map[string]int64{
		"add_num": int64(10),
	})

	inc1 := metrics.Get("inc_num")
	inc1.Add(4)

	checkFileWriter(t, file.Name(), lineSep, map[string]int64{
		"add_num": int64(10),
		"inc_num": int64(4),
	})
}

func TestMetricsSetFormatter(t *testing.T) {
	t.Parallel()

	file := newTempFile(t)
	fileName := file.Name()
	defer removeTempFile(t, file)

	metrics := New()
	metrics.SetOutput(file)
	metrics.SetFormatter(NewFormatter("\n"))

	c := metrics.Get("test_counter")
	c.Add(10)

	require.Nil(t, metrics.Write())

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
	t.Parallel()

	metrics := New()
	metrics.SetFormatter(NewFormatter("\n"))
	assert.NotNil(t, metrics.Formatter())
}

func TestMetricsDefault(t *testing.T) {
	t.Parallel()

	file := newTempFile(t)

	SetOutput(file)
	SetFormatter(NewFormatter("\n"))

	require.NotNil(t, Default.formatter)

	c := Get("default_metrics_counter")
	c.Add(10)

	require.Nil(t, Write())
	removeTempFile(t, file)

	SetPanicHandler(new(mockPanicHandler))
	assert.NotNil(t, Default.panicHandler)

	file = newTempFile(t)
	defer removeTempFile(t, file)

	StartFileWriter(FileWriterParams{
		FilePath:       file.Name(),
		UpdateInterval: time.Millisecond * 100,
	})
	StopFileWriter()

	counter := Get("default_metrics_counter")
	require.NotNil(t, counter)
	require.Equal(t, int64(10), counter.Get())

	prefixMetrics := WithPrefix("prefix.%s.", "errors")
	c2 := prefixMetrics.Get("data")
	assert.True(t, c2 == Get("prefix.errors.data"))
}

func TestMetricsSetPanicHandler(t *testing.T) {
	t.Parallel()

	metrics := New()
	handler := new(mockPanicHandler)
	metrics.SetPanicHandler(handler)
	assert.Equal(t, handler, metrics.getPanicHandler())
}

func TestMetricsGetTwice(t *testing.T) {
	t.Parallel()

	metrics := New()
	c := metrics.Get("new_counter")
	require.NotNil(t, c)
	require.Equal(t, int64(0), c.Get())
	c.Set(11)

	c2 := metrics.Get("new_counter")
	require.NotNil(t, c2)
	require.Equal(t, int64(11), c2.Get())
	assert.True(t, c == c2)
}

func TestMetricsGetJSON(t *testing.T) {
	t.Parallel()

	metrics := New()

	counter1 := metrics.Get("counter1")
	counter1.Set(10)

	counter2 := metrics.Get("counter2")
	counter2.Set(42)

	b := metrics.GetJSON(func(key string) bool {
		if key == "counter1" || key == "counter2" {
			return true
		}
		return false
	})
	assert.JSONEq(t, `{"counter1": 10, "counter2": 42}`, string(b))

	b = metrics.GetJSON(func(key string) bool {
		if key == "counter2" {
			return true
		}
		return false
	})
	assert.JSONEq(t, `{"counter2": 42}`, string(b))

	b = metrics.GetJSON(func(string) bool {
		return false
	})
	assert.JSONEq(t, `{}`, string(b))
}

func TestMetricsDefaultGetJSON(t *testing.T) {
	t.Parallel()

	counter1 := Get("counter1")
	counter1.Set(10)

	counter2 := Get("counter2")
	counter2.Set(42)

	b := GetJSON(func(key string) bool {
		if key == "counter1" || key == "counter2" {
			return true
		}
		return false
	})
	assert.JSONEq(t, `{"counter1": 10, "counter2": 42}`, string(b))

	b = GetJSON(func(key string) bool {
		if key == "counter2" {
			return true
		}
		return false
	})
	assert.JSONEq(t, `{"counter2": 42}`, string(b))

	b = GetJSON(func(key string) bool {
		return false
	})
	assert.JSONEq(t, `{}`, string(b))
}

func TestMetricsSetRootPrefix(t *testing.T) {
	t.Parallel()

	file := newTempFile(t)
	require.Nil(t, file.Close())
	defer os.Remove(file.Name())

	prefix := "test."
	metrics := New()
	metrics.SetRootPrefix(prefix)
	lineSep := "\n"

	inc := metrics.Get("add_num")
	inc.Add(10)

	metrics.StartFileWriter(FileWriterParams{
		FilePath:       file.Name(),
		UpdateInterval: time.Millisecond * 100,
	})
	defer metrics.StopFileWriter()

	checkFileWriter(t, file.Name(), lineSep, map[string]int64{
		prefix + "add_num": int64(10),
	})
}

type mockPanicHandler struct{}

func (h *mockPanicHandler) Handle(err error) {
	fmt.Fprintf(os.Stderr, "failed to write metrics file, %v\n", err)
}

func newTempFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", t.Name())
	require.Nil(t, err)
	return file
}

func removeTempFile(t *testing.T, f *os.File) {
	require.Nil(t, f.Close())
	require.Nil(t, os.Remove(f.Name()))
}

func checkFileWriter(t *testing.T, fileName, lineSep string, counters map[string]int64) {
	ch := time.After(time.Minute)
	var updateNum int

	check := func() bool {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			if !os.IsNotExist(err) {
				require.FailNow(t, err.Error())
			}
			return false
		}

		metrics := strings.TrimSuffix(string(data), lineSep)
		metricsData := strings.Split(metrics, lineSep)

		for _, metricLine := range metricsData {
			if metricLine != "" {
				updateNum++

				metricLine = strings.Replace(metricLine, " ", "", -1)
				counter := strings.Split(metricLine, "=")
				require.Len(t, counter, 2)

				key := counter[0]
				actualVal, err := strconv.ParseInt(counter[1], 10, 64)
				require.Nil(t, err)

				expectedVal, ok := counters[key]
				require.True(t, ok)
				require.Equal(t, expectedVal, actualVal)

				if updateNum == len(counters) {
					return true
				}
			}
		}

		return false
	}

	for {
		select {
		case <-ch:
			require.FailNow(t, "counters wasn't updated")
			return
		default:
			if check() {
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}
