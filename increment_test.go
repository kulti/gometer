package gometer

import (
	"strings"
	"testing"

	"os"

	"io/ioutil"

	"strconv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Create(fileName)
	require.Nil(t, err)
	return file
}

type testCounterParams struct {
	metricName     string
	operationID    string
	fileName       string
	operationCount int
	initialValue   int64
	expectedValue  int
}

func testCounter(t *testing.T, p testCounterParams) {
	var c Incrementor
	if p.operationID == "inc" {
		c = NewIncrementor(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Inc()
		}
	} else if p.operationID == "add" {
		c = NewIncrementor(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Add(p.initialValue)
		}
	} else if p.operationID == "set" {
		c := NewCounter(p.metricName)
		for i := 0; i < p.operationCount; i++ {
			c.Set(p.initialValue)
		}
	} else if p.operationID == "dec" {
		c := NewCounter(p.metricName)
		c.Set(p.initialValue)
		for i := 0; i < p.operationCount; i++ {
			c.Dec()
		}
	}

	// write all existing metrics to output
	err := Write()
	require.Nil(t, err)

	data, err := ioutil.ReadFile(p.fileName)
	require.Nil(t, err)

	lines := strings.Split(string(data), "\n")
	require.True(t, len(lines) > 1)

	for _, l := range lines {
		counterData := strings.Split(l, " = ")
		require.True(t, len(counterData) == 2)
		if counterData[0] == p.metricName {
			// check the counter value
			val, err := strconv.Atoi(counterData[1])
			require.Nil(t, err)
			assert.Equal(t, p.expectedValue, val)
			return
		}
	}
}

func closeAndRemoveTestFile(t *testing.T, f *os.File) {
	err := f.Close()
	require.Nil(t, err)
	err = os.Remove(f.Name())
	require.Nil(t, err)
}

func TestInc(t *testing.T) {
	file := newTestFile(t, "test_inc")
	defer closeAndRemoveTestFile(t, file)

	SetOutput(file)
	SetFormat("%v = %v\n")

	testCounter(t, testCounterParams{
		metricName:     "simple_counter1",
		operationCount: 10,
		operationID:    "inc",
		fileName:       file.Name(),
		expectedValue:  10,
	})
}

func TestAdd(t *testing.T) {
	file := newTestFile(t, "test_add")
	defer closeAndRemoveTestFile(t, file)

	SetOutput(file)
	SetFormat("%v = %v\n")

	testCounter(t, testCounterParams{
		metricName:     "simple_adder",
		operationCount: 2,
		operationID:    "add",
		fileName:       file.Name(),
		expectedValue:  4,
		initialValue:   2,
	})
}
