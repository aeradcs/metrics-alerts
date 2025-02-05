package client

import (
	"runtime"
)

type ExtendedMemStats struct {
	runtime.MemStats
	PollCount   int64
	RandomValue float64
}
