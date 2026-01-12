package utils

import (
	"math/rand/v2"
	"time"
)

// RandomDuration generates a random duration between 1 minute and 5 minutes
func RandomDuration() time.Duration {
	minDuration := time.Minute
	maxDuration := 5 * time.Minute
	
	// Generate random duration in nanoseconds
	diff := maxDuration - minDuration
	randomNanos := rand.Int64N(int64(diff))
	
	return minDuration + time.Duration(randomNanos)
}
