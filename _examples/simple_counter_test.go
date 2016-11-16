package example

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/dshil/gometer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type checkMetricParams struct {
	name        string
	value       int64
	metricsData []string
}

func checkMetric(t *testing.T, p checkMetricParams) {
	for _, m := range p.metricsData {
		metric := strings.Split(m, " = ")
		require.Equal(t, 2, len(metric))

		metricName := metric[0]
		if metricName == p.name {
			// check value
			metricVal := metric[1]
			val, err := strconv.Atoi(metricVal)
			require.Nil(t, err)
			assert.Equal(t, p.value, int64(val))
		}
	}
}

func TestSimpleCounter(t *testing.T) {
	// init test file where to dump all metrics.
	fileName := "test_file"
	file, err := os.Create(fileName)
	require.Nil(t, err)
	defer file.Close()
	defer os.Remove(fileName)

	// make some preparation for standard gometer.
	gometer.SetOutput(file)

	// choose a format of metric representation.
	// e.g metric_name = metric_value.
	gometer.SetFormat("%v = %v")

	// each metric line will be separated by \n.
	gometer.SetSeparator("\n")

	// require to call gometer.Write() method manually
	// because update interval equals to 0.
	gometer.SetUpdateInterval(0)

	// init simple counter and increment it 10 times.
	inc := gometer.NewIncrementor("number_incrementor")
	for i := 0; i < 10; i++ {
		inc.Inc()
	}
	assert.Equal(t, int64(10), inc.Value())

	dec := gometer.NewCounter("number_decrementor")
	dec.Set(5)
	dec.Dec()
	assert.Equal(t, int64(4), dec.Value())

	// write all metrics to file.
	err = gometer.Write()
	require.Nil(t, err)

	// need to check if file contains the right values for metrics.
	data, err := ioutil.ReadFile(fileName)
	require.Nil(t, err)

	// metrics are splitted using \n separator.
	// need to trim separator from last line of the file.
	metrics := strings.TrimSuffix(string(data), gometer.Separator())
	metricsData := strings.Split(metrics, gometer.Separator())

	// we have only 2 metrics in file.
	require.Equal(t, 2, len(metricsData))

	// check the corresponding names and values for metrics.
	checkMetric(t, checkMetricParams{
		name:        "number_incrementor",
		value:       inc.Value(),
		metricsData: metricsData,
	})
	checkMetric(t, checkMetricParams{
		name:        "number_decrementor",
		value:       dec.Value(),
		metricsData: metricsData,
	})
}
