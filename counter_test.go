package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterAdd(t *testing.T) {
	c := Counter{}
	c.Add(10)
	assert.Equal(t, int64(10), c.Get())
}

func TestCounterSet(t *testing.T) {
	c := Counter{}
	c.Set(10)
	assert.Equal(t, int64(10), c.Get())

	c.Set(-10)
	assert.Equal(t, int64(-10), c.Get())
}
