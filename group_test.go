package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountersGroupAdd(t *testing.T) {
	metrics := New()
	group := metrics.Group("foo.")

	counter := DefaultCounter{}
	counter.Add(10)
	group.Add("bar", &counter)

	got := group.Get("bar")
	require.NotNil(t, got)
	assert.Equal(t, int64(10), got.Get())
}

func TestCountersGroupCounters(t *testing.T) {
	metrics := New()

	group := metrics.Group("%s.", "foo")
	require.NotNil(t, group)

	bazCounter := DefaultCounter{}
	bazCounter.Set(10)

	group.Add("baz", &bazCounter)

	groupCounters := group.Counters()
	require.Len(t, groupCounters, 1)

	gotBazCounter := groupCounters["foo.baz"]
	require.NotNil(t, gotBazCounter)

	assert.Equal(t, int64(10), gotBazCounter.Get())
}
