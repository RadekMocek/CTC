package main

import (
	"fmt"
	"sync"
)

var mainWG sync.WaitGroup

// Spawn `nCars` of cars in intervals between `minTime` and `maxTime`; add them to sharedQueue `sq` if there is space; if not, wait (block)
func carSpawner(nCars int, minTime int, maxTime int, sq *sharedQueue) {
	for i := 0; i < nCars; i++ {
		randSleepExcl(minTime, maxTime)
		fmt.Println("Car with id=", i, "wants to enter the shared queue.")
		//sq.queue <- &car{i, randFuelType()}
		sq.queue <- &car{i, gas}
	}
}

// Get index of stand or register where it is best for the car to go (shortest queue)
func getBestIndex(specificStands []*standOrRegister) int {
	var tempValue int
	// If all of the stands are full, driver chooses the first one; they dont know which stand will be free the first
	bestIndex := 0
	// ISSUE(?): We're getting the len of channel `queue`; at the same time `standFillAndDistributeToRegisters` could read from it (and decrease the len)
	bestValue := len(specificStands[0].queue)
	if specificStands[0].isUsed {
		bestValue++
	}
	for i := 1; i < len(specificStands); i++ {
		tempValue = len(specificStands[i].queue)
		if specificStands[i].isUsed {
			tempValue++
		}
		if tempValue < bestValue {
			bestValue = tempValue
			bestIndex = i
		}
	}
	return bestIndex
}

// Used as a single goroutine; distributes cars from sharedQueue `sq` to stands accroding to their fuel type
func (sq *sharedQueue) distributeToStands(gasStands []*standOrRegister) {
	var bestIndex int
	var specificStands []*standOrRegister

	for {
		c := <-sq.queue
		fmt.Println("Car", c, "enters the shared queue.")
		switch c.fuelType {
		case gas:
			specificStands = gasStands
		}
		bestIndex = getBestIndex(specificStands)
		fmt.Println("Car", c, "CHOOSES stand", specificStands[bestIndex])
		specificStands[bestIndex].queue <- c
		fmt.Println("Car", c, "ENTERS stand", specificStands[bestIndex])
	}
}

// Used as goroutine for every stand; simulates refueling and after that sends cars to cash registers
func (sor *standOrRegister) standFillAndDistributeToRegisters(registers []*standOrRegister) {
	var bestIndex int

	for {
		c := <-sor.queue
		sor.isUsed = true
		fmt.Println("Car", c, "REFUELING in stand", sor)
		randSleepExcl(sor.minTime, sor.maxTime)

		bestIndex = getBestIndex(registers)
		fmt.Println("Car", c, "CHOOSES register with id=", bestIndex)
		registers[bestIndex].queue <- c
		fmt.Println("Car", c, "ENTERS register with id=", bestIndex)

		sor.isUsed = false
	}
}

// Used as goroutine for every cash register; simulates paying and then makes cars leave the gas station
func (sor *standOrRegister) registerCashoutAndLeave() {
	for {
		c := <-sor.queue
		sor.isUsed = true
		fmt.Println("Car", c, "PAYING in register", sor)
		randSleepExcl(sor.minTime, sor.maxTime)
		mainWG.Done()
		fmt.Println("Car", c, "LEAVES register", sor)
		sor.isUsed = false
	}
}

func main() {
	fmt.Println("Hello gas station!")

	// Read data from config.yaml to the conf variable
	conf, err := readConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Retrieve some repeatedly used values from the conf variable
	carsConf := conf.Cars
	nCars := carsConf.Count
	gasConf := conf.Stations.Gas
	nGasStands := gasConf.Count
	registerConf := conf.Registers
	nRegisters := registerConf.Count

	// Simulation
	// Make cash registers and start their goroutines that accept cars from stands and makes them leave the gas station after their payment is done
	registers := make([]*standOrRegister, nRegisters)
	for i := 0; i < nRegisters; i++ {
		registers[i] = &standOrRegister{i, false, make(chan *car, registerConf.QueueLengthMax), registerConf.HandleTimeMin, registerConf.HandleTimeMax}
		go registers[i].registerCashoutAndLeave()
	}

	// Make gasStands and start their goroutines that accept cars from sharedQueue, refuels them, and sends them to cash registers
	gasStands := make([]*standOrRegister, nGasStands)
	for i := 0; i < nGasStands; i++ {
		gasStands[i] = &standOrRegister{(gas+1)*10 + i, false, make(chan *car, gasConf.QueueLengthMax), gasConf.ServeTimeMin, gasConf.ServeTimeMax}
		go gasStands[i].standFillAndDistributeToRegisters(registers)
	}

	// Make a sharedQueue and start a goroutine that sends cars to specific queues according to their fuel type
	sq := sharedQueue{make(chan *car, carsConf.SharedQueueLengthMax)}
	go sq.distributeToStands(gasStands)

	// Start spawning cars
	mainWG.Add(nCars)
	go carSpawner(nCars, carsConf.ArrivalTimeMin, carsConf.ArrivalTimeMax, &sq)
	mainWG.Wait()

	// Amen
	fmt.Println("Simulation done.")
}
