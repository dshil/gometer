package gometer

import "testing"

func TestSet(t *testing.T) {
	file := newTestFile(t, "test_set")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)
	SetFormat("%v = %v\n")

	testCounter(t, testCounterParams{
		metricName:     "test_set_positive_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  100,
		operationValue: 100,
	})

	testCounter(t, testCounterParams{
		metricName:     "test_set_negative_val",
		operationCount: 1,
		operationID:    "set",
		fileName:       file.Name(),
		expectedValue:  -19,
		operationValue: -19,
	})
}

func TestDec(t *testing.T) {
	file := newTestFile(t, "test_dec")
	defer closeAndRemoveTestFile(t, file)
	SetOutput(file)
	SetFormat("%v = %v\n")

	testCounter(t, testCounterParams{
		metricName:     "test_decrement_counter",
		operationID:    "dec",
		fileName:       file.Name(),
		expectedValue:  2,
		operationValue: 4,
		operationCount: 2,
	})
}
