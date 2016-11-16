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
	m := New()

	inc := m.NewIncrementor("add_num")
	inc.Add(10)

	m.WriteAtFile(fileName)
	testWriteAtFile(t, fileName, m.Separator(), 1)

	inc1 := m.NewIncrementor("inc_num")
	inc1.Inc()

	m.WriteAtFile(fileName)
	testWriteAtFile(t, fileName, m.Separator(), 2)
}

func testWriteAtFile(t *testing.T, fileName, separator string, expMetricsCnt int) {
	data, err := ioutil.ReadFile(fileName)
	require.Nil(t, err)
	defer os.Remove(fileName)
	metrics := strings.TrimSuffix(string(data), separator)
	metricsData := strings.Split(metrics, separator)
	require.Equal(t, expMetricsCnt, len(metricsData))
}
