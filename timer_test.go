package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_time(t *testing.T) {
	time := startTimer()
	assert.NotZero(t, time.startTime)
	assert.Zero(t, time.finishTime)
	diff := time.diff()
	assert.NotZero(t, diff)
	time.stop()
	assert.NotZero(t, time.finishTime)
}
