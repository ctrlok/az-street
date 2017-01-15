package main

import "time"

type timer struct {
	startTime  int64
	finishTime int64
}

func (t *timer) stop() int64 {
	t.finishTime = t.diff()
	return t.finishTime
}

func (t *timer) diff() int64 {
	return (time.Now().UnixNano() - t.startTime)
}

func startTimer() *timer {
	t := timer{}
	t.startTime = time.Now().UnixNano()
	return &t
}
