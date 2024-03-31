package main

import (
	"math/rand"
	"time"
)

const durationScaler time.Duration = time.Millisecond

// Sleep for a random duration between `minTime` and `maxTime`; both exclusive; multiplied by const durationScaler
func randSleepExcl(minTime int, maxTime int) {
	time.Sleep(time.Duration(rand.Intn(maxTime+1-minTime)+minTime) * durationScaler)
}

func randFuelType() int {
	return rand.Intn(4)
}
