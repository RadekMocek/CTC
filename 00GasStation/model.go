package main

const (
	gas      = 0
	diesel   = 1
	lpg      = 2
	electric = 3
)

type car struct {
	id       int
	fuelType int
}

type sharedQueue struct {
	queue chan *car
}

type standOrRegister struct {
	id int
	//fuelType int
	isUsed  bool
	queue   chan *car
	minTime int
	maxTime int
}
