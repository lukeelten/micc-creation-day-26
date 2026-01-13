package utils

import (
	"context"
	"math/rand/v2"
	"time"
)

// RandomDuration generates a random duration between 1 minute and 5 minutes
func RandomDuration() time.Duration {
	return RandomDurationLimit(5 * time.Minute)
}

// RandomDuration generates a random duration between 1 minute and 5 minutes
func RandomDurationLimit(maxDuration time.Duration) time.Duration {
	minDuration := 5 * time.Second

	// Generate random duration in nanoseconds
	diff := maxDuration - minDuration
	randomNanos := rand.Int64N(int64(diff))

	return minDuration + time.Duration(randomNanos)
}

func SimulateWork(ctx context.Context, duration time.Duration) {
	// Allocate ~64 MiB: 64 * 1024 * 1024 bytes.
	size := 64 * 1024 * 1024
	buf := make([]byte, size)

	// Touch the memory so the OS commits the pages.
	for i := range size {
		if i%4096 == 0 {
			buf[i] = byte(i)
		}
	}

	deadline := time.Now().Add(duration)

	// Minor CPU cycles while waiting.
	var acc uint64
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return
		default:
			for i := range 1000 {
				acc += uint64(buf[(i*4096)%size])
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Prevent compiler from optimizing away acc/buf.
	_ = acc
}
