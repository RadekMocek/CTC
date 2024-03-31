package main

import "math/rand"

// Returns random int between min and max; both exclusive
func randRangeExcl(min int, max int) int {
	return rand.Intn(max+1-min) + min
}
