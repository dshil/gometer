package gometer

import (
	"strings"
	"testing"

	"os"

	"io/ioutil"

	"strconv"

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
	var c Incrementor
	if p.operationID == "inc" {
		c = NewIncrementor(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Inc()
		}
	} else if p.operationID == "add" {
		c = NewIncrementor(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Add(p.initialValue)
		}
	} else if p.operationID == "set" {
		c := NewCounter(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Set(p.initialValue)
		}
	} else if p.operationID == "dec" {
		c := NewCounter(p.metricName)
		c.Set(p.initialValue)
		for i := 0; i < p.operationCount; i++ {
			c.Dec()
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

func TestInc(t *testing.T) {
	file := newTestFile(t, "test_inc")
	defer closeAndRemoveTestFile(t, file)

	SetOutput(file)
	SetFormat("%v = %v")
	SetSeparator("\n")

	testCounter(t, testCounterParams{
		metricName:     "simple_counter1",
		operationCount: 10,
		operationID:    "inc",
		fileName:       file.Name(),
		expectedValue:  10,
	})
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
