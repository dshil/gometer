package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterAdd(t *testing.T) {
	metrics := New()
	c := metrics.NewCounter("add")
	c.Add(10)
	assert.Equal(t, int64(10), c.Get())
}

func TestCounterSet(t *testing.T) {
	metrics := New()
	c := metrics.NewCounter("set")
	c.Set(10)
	assert.Equal(t, int64(10), c.Get())

	c.Set(-10)
	assert.Equal(t, int64(-10), c.Get())
}
