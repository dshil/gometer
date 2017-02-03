package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupCounterWithPrefix(t *testing.T) {
	metrics := New()

	group := metrics.Group().WithPrefix("foo.%s.baz", "bar")

	require.NotNil(t, group)
	assert.Equal(t, "foo.bar.baz", group.Prefix())
}

func TestGroupCounterAdd(t *testing.T) {
	metrics := New()
	group := metrics.Group().WithPrefix("foo")

	counter := Counter{}
	counter.Add(10)
	group.Add("bar", &counter)

	got := group.Get("bar")
	require.NotNil(t, got)
	assert.Equal(t, int64(10), got.Get())
}

func TestGroupCounterWithSeparator(t *testing.T) {
	metrics := New()
	group := metrics.Group().WithPrefix("foo.%s", "bar").WithSeparator(".")
	require.NotNil(t, group)

	assert.Equal(t, ".", group.Separator())
}

func TestGroupCounterCounters(t *testing.T) {
	metrics := New()
	group := metrics.Group().WithPrefix("foo.%s", "bar").WithSeparator(".")
	require.NotNil(t, group)

	bazCounter := Counter{}
	bazCounter.Set(10)

	group.Add("baz", &bazCounter)

	groupCounters := group.Counters()
	require.Len(t, groupCounters, 1)

	gotBazCounter := groupCounters["foo.bar.baz"]
	require.NotNil(t, gotBazCounter)

	assert.Equal(t, int64(10), gotBazCounter.Get())
}
