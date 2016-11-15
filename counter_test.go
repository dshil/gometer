package gometer

import "testing"

func TestSet(t *testing.T) {
	file := newTestFile(t, "test_set_positive")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)
	SetFormat("%v = %v")

	testCounter(t, testCounterParams{
		metricName:     "test_set_positive_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  100,
		initialValue:   100,
	})

	file = newTestFile(t, "test_set_negative")
	SetOutput(file)

	defer closeAndRemoveTestFile(t, file)
	testCounter(t, testCounterParams{
		metricName:     "test_set_negative_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  -19,
		initialValue:   -19,
	})
}

func TestDec(t *testing.T) {
	file := newTestFile(t, "test_dec")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)
	SetFormat("%v = %v")

	testCounter(t, testCounterParams{
		metricName:     "test_decrement_counter",
		operationID:    "dec",
		fileName:       file.Name(),
		expectedValue:  2,
		initialValue:   4,
		operationCount: 2,
	})
}
