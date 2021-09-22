package main

import (
	"time"

	"go.uber.org/zap"
)

type constantClock time.Time

func (c constantClock) Now() time.Time {
	return time.Time(c)
}
func (c constantClock) NewTicker(duration time.Duration) *time.Ticker {
	return &time.Ticker{}
}

func ExampleNewMultiOutputLogger() {
	date := time.Date(2021, 1, 23, 10, 15, 13, 0, time.UTC)
	clock := constantClock(date)
	logger := NewMultiOutputLogger()
	logger = logger.WithOptions(zap.WithClock(clock))
	logger.Info("info message")

	// Output:
	// 2021-01-23T10:15:13.000Z	INFO	info message
}
