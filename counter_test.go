package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterAdd(t *testing.T) {
	c := DefaultCounter{}
	c.Add(10)
	assert.Equal(t, int64(10), c.Get())
}

func TestCounterSet(t *testing.T) {
	c := DefaultCounter{}
	c.Set(10)
	assert.Equal(t, int64(10), c.Get())

	c.Set(-10)
	assert.Equal(t, int64(-10), c.Get())
}

func TestCounterAddAndGet(t *testing.T)  {
	c := DefaultCounter{}
	v := c.AddAndGet(10)
	assert.Equal(t, int64(10), v)
	assert.Equal(t, int64(10), c.Get())

	v = c.AddAndGet(5)
	assert.Equal(t, int64(15), v)
	assert.Equal(t, int64(15), c.Get())
}
