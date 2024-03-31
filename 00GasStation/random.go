package main

import (
	"math/rand"
	"time"
)

const durationScaler time.Duration = time.Second

// Returns random Duration between `min` and `max`; both exclusive
func randDurExcl(min int, max int) time.Duration {
	return time.Duration(rand.Intn(max+1-min) + min)
}

// Sleep for a random duration between `minTime` and `maxTime`; both exclusive; multiplied by const durationScaler
func randSleepExcl(minTime int, maxTime int) {
	time.Sleep(randDurExcl(minTime, maxTime) * durationScaler)
}

func randFuelType() int {
	return rand.Intn(4)
}
