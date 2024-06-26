package main

import (
	"time"
)

const (
	gas      = 0
	diesel   = 1
	lpg      = 2
	electric = 3
)

type car struct {
	id       int
	fuelType int
	// Stats
	waitForSharedQueueStarted time.Time
	waitForStandStarted       time.Time
	waitForRegisterStarted    time.Time
	departureStarted          time.Time
}

type sharedQueue struct {
	queue chan car
}

type standOrRegister struct {
	id      int
	isUsed  bool
	queue   chan car
	minTime int
	maxTime int
}

type standConfigRepresentation struct {
	Count          int
	ServeTimeMin   int
	ServeTimeMax   int
	QueueLengthMax int
}
