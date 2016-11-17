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

	inc := m.NewCounter("add_num")
	inc.Add(10)

	stopper := m.WriteToFile(fileName, time.Second*10, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:     fileName,
		separator:    m.Separator(),
		expMetricCnt: 1,
		waitDur:      time.Second * 1,
		stopper:      stopper,
	})

	inc1 := m.NewCounter("inc_num")
	inc1.Add(4)

	stopper = m.WriteToFile(fileName, time.Second*10, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:     fileName,
		separator:    m.Separator(),
		expMetricCnt: 2,
		waitDur:      time.Second * 1,
		stopper:      stopper,
	})
}

type testWriteAtFileParams struct {
	fileName     string
	separator    string
	expMetricCnt int
	waitDur      time.Duration
	stopper      *stopper
}

func testWriteAtFile(t *testing.T, p testWriteAtFileParams) {
	time.Sleep(p.waitDur)
	p.stopper.Stop()

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	err = os.Remove(p.fileName)
	require.Nil(t, err)

	metrics := strings.TrimSuffix(string(data), p.separator)
	metricsData := strings.Split(metrics, p.separator)
	require.Equal(t, p.expMetricCnt, len(metricsData))
}
