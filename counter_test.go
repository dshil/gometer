package gometer

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Create(fileName)
	require.Nil(t, err)
	return file
}

type testCounterParams struct {
	counterName string
	operationID string

	fileName string

	operationCount int
	metricsCount   int

	initialValue  int64
	expectedValue int
}

func testCounter(t *testing.T, p testCounterParams) {
	if p.operationID == "add" {
		c := NewCounter(p.counterName)
		for i := 0; i < p.operationCount; i++ {
			c.Add(p.initialValue)
		}
	} else if p.operationID == "set" {
		c := NewCounter(p.counterName)
		for i := 0; i < p.operationCount; i++ {
			c.Set(p.initialValue)
		}
	}

	// write all existing metrics to output
	err := Write()
	require.Nil(t, err)

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	metrics := strings.TrimSuffix(string(data), Formatter().LineSeparator)
	metricsData := strings.Split(metrics, Formatter().LineSeparator)

	var reqMetricLen bool
	if p.metricsCount == 0 {
		reqMetricLen = len(metricsData) >= 1
	} else {
		reqMetricLen = (len(metricsData) == p.metricsCount)
	}
	require.True(t, reqMetricLen)

	for _, l := range metricsData {
		metric := strings.Split(l, " = ")
		require.True(t, len(metric) == 2)
		metricName := metric[0]
		if metricName == p.counterName {
			// check the counter value
			metricVal := metric[1]
			val, err := strconv.Atoi(metricVal)
			require.Nil(t, err)
			assert.Equal(t, p.expectedValue, val)
			return
		}
	}
}

func closeAndRemoveTestFile(t *testing.T, f *os.File) {
	err := f.Close()
	require.Nil(t, err)
	err = os.Remove(f.Name())
	require.Nil(t, err)
}

func TestAdd(t *testing.T) {
	file := newTestFile(t, "test_add")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)

	testCounter(t, testCounterParams{
		counterName:    "simple_adder",
		operationCount: 2,
		operationID:    "add",
		fileName:       file.Name(),
		expectedValue:  4,
		initialValue:   2,
	})
}

func TestSet(t *testing.T) {
	file := newTestFile(t, "test_set_positive")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)

	testCounter(t, testCounterParams{
		counterName:    "test_set_positive_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  100,
		initialValue:   100,
	})

	testCounter(t, testCounterParams{
		counterName:    "test_set_negative_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  -19,
		initialValue:   -19,
	})
}

func TestGet(t *testing.T) {
	m := New()
	counter := m.NewCounter("get_metric")
	counter.Add(10)
	assert.Equal(t, int64(10), counter.Get())
}
