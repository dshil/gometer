package gometer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWriteAtFile(t *testing.T) {
	fileName := "test_write_at_file"
	m := New()

	inc := m.NewIncrementor("add_num")
	inc.Add(10)

	stopper := m.WriteAtFile(fileName, time.Second*2, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:     fileName,
		separator:    m.Separator(),
		expMetricCnt: 1,
		waitDur:      1,
		stopper:      stopper,
	})

	inc1 := m.NewIncrementor("inc_num")
	inc1.Inc()

	m.WriteAtFile(fileName, time.Second, true)
	testWriteAtFile(t, testWriteAtFileParams{
		fileName:     fileName,
		separator:    m.Separator(),
		expMetricCnt: 2,
		waitDur:      1,
		stopper:      stopper,
	})
}

type testWriteAtFileParams struct {
	fileName     string
	separator    string
	expMetricCnt int
	waitDur      time.Duration
	stopper      *Stopper
}

func testWriteAtFile(t *testing.T, p testWriteAtFileParams) {
	time.Sleep(time.Second * p.waitDur)
	p.stopper.Stop()

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	err = os.Remove(p.fileName)
	require.Nil(t, err)

	metrics := strings.TrimSuffix(string(data), p.separator)
	metricsData := strings.Split(metrics, p.separator)
	require.Equal(t, p.expMetricCnt, len(metricsData))
}
