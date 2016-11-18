package gometer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWriteToFile(t *testing.T) {
	fileName := "test_write_to_file"
	m := New()
	m.SetFormatter(NewDefaultFormatter())

	inc := m.NewCounter("add_num")
	inc.Add(10)

	stopper := m.WriteToFile(fileName, time.Second*10, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:      fileName,
		lineSeparator: m.Formatter().LineSeparator,
		expMetricCnt:  1,
		waitDur:       time.Second * 1,
		stopper:       stopper,
	})

	inc1 := m.NewCounter("inc_num")
	inc1.Add(4)

	stopper = m.WriteToFile(fileName, time.Second*10, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:      fileName,
		lineSeparator: m.Formatter().LineSeparator,
		expMetricCnt:  2,
		waitDur:       time.Second * 1,
		stopper:       stopper,
	})
}

type testWriteAtFileParams struct {
	fileName      string
	lineSeparator string

	expMetricCnt int

	waitDur time.Duration
	stopper *Stopper
}

func testWriteAtFile(t *testing.T, p testWriteAtFileParams) {
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
