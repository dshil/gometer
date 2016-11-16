package gometer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteAtFile(t *testing.T) {
	fileName := "test_write_at_file"
	SetSeparator("\n")
	SetFormat("%v = %v")

	inc := NewIncrementor("add_num")
	inc.Add(10)

	WriteAtFile(fileName)
	testWriteAtFile(t, fileName, 1)

	inc1 := NewIncrementor("inc_num")
	inc1.Inc()

	WriteAtFile(fileName)
	testWriteAtFile(t, fileName, 2)
}

func testWriteAtFile(t *testing.T, fileName string, expMetricsCnt int) {
	data, err := ioutil.ReadFile(fileName)
	require.Nil(t, err)
	defer os.Remove(fileName)
	metrics := strings.TrimSuffix(string(data), Separator())
	metricsData := strings.Split(metrics, Separator())
	require.Equal(t, expMetricsCnt, len(metricsData))
}

func ExampleWriteToStdout() {
	SetFormat("%v = %v")
	SetOutput(os.Stdout)
	SetSeparator("\n")
	c := NewCounter("err_cnt")
	c.Inc()
	Write()
	// Output: err_cnt = 1
}
