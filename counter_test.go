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
	metricName     string
	operationID    string
	fileName       string
	operationCount int
	initialValue   int64
	expectedValue  int
	metricNumber   int
}

func testCounter(t *testing.T, p testCounterParams) {
	if p.operationID == "add" {
		c := NewCounter(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Add(p.initialValue)
		}
	} else if p.operationID == "set" {
		c := NewCounter(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Set(p.initialValue)
		}
	}

	// write all existing metrics to output
	err := Write()
	require.Nil(t, err)

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	metrics := strings.TrimSuffix(string(data), Separator())
	metricsData := strings.Split(metrics, Separator())

	var reqMetricLen bool
	if p.metricNumber == 0 {
		reqMetricLen = len(metricsData) >= 1
	} else {
		reqMetricLen = (len(metricsData) == p.metricNumber)
	}
	require.True(t, reqMetricLen)

	for _, l := range metricsData {
		metric := strings.Split(l, " = ")
		require.True(t, len(metric) == 2)
		metricName := metric[0]
		if metricName == p.metricName {
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
	SetFormat("%v = %v")
	SetSeparator("\n")

	testCounter(t, testCounterParams{
		metricName:     "simple_adder",
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
	SetSeparator("\n")
	SetFormat("%v = %v")

	testCounter(t, testCounterParams{
		metricName:     "test_set_positive_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  100,
		initialValue:   100,
	})

	testCounter(t, testCounterParams{
		metricName:     "test_set_negative_val",
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
