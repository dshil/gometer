package gometer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefixMetricsGet(t *testing.T) {
	originalMetrics := New()
	prefixMetrics := originalMetrics.WithPrefix("data.")

	counter1 := originalMetrics.Get("data.first")
	counter2 := originalMetrics.Get("data.second")

	require.True(t, counter1 == prefixMetrics.Get("first"))
	assert.True(t, counter2 == prefixMetrics.Get("second"))
}

func TestPrefixMetricsGetFormatted(t *testing.T) {
	originalMetrics := New()
	prefixMetrics := originalMetrics.WithPrefix("data.%s.%s.", "test", "errors")

	c := prefixMetrics.Get("counter")
	assert.True(t, c == originalMetrics.Get("data.test.errors.counter"))
}

func TestPrefixMetricsTwice(t *testing.T) {
	originalMetrics := New()
	prefixMetrics1 := originalMetrics.WithPrefix("prefix1.").WithPrefix("prefix2.%s.", "errors")

	c := prefixMetrics1.Get("counter")
	assert.True(t, c == originalMetrics.Get("prefix1.prefix2.errors.counter"))
}
