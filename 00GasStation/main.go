package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var globalWG sync.WaitGroup

// Spawn `nCars` of cars in intervals between `minTime` and `maxTime`; add them to sharedQueue `sq` if there is space; if not, wait (block)
func carSpawner(nCars int, minTime int, maxTime int, sq *sharedQueue) {
	var c car
	for i := 0; i < nCars; i++ {
		randSleepExcl(minTime, maxTime)
		c = car{}
		c.id = i
		c.fuelType = randFuelType()
		// Car `c` was created and wants to enter the sharedQueue
		c.waitForSharedQueueStarted = time.Now()
		sq.queue <- c
		// Car `c` has entered the back of the sharedQueue
	}
}

// Get index of stand or register where it is best for the car to go (shortest queue)
func getBestIndex(specificStands []*standOrRegister) int {
	var tempValue int
	// If all of the stands are full, driver chooses the first one; they dont know which stand will be free the first
	bestIndex := 0
	// ISSUE: We're getting the len of channel `queue`; at the same time `standFillAndDistributeToRegisters` could read from it (and decrease the len)
	// https://stackoverflow.com/a/42321398
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
func (sq *sharedQueue) distributeToStands(allStands map[int][]*standOrRegister) {
	var bestIndex int
	var c car
	var specificStands []*standOrRegister
	for {
		c = <-sq.queue
		// Car `c` is at the front of the sharedQueue
		specificStands = allStands[c.fuelType] // Look for stand that corresponds to car's fuel type
		bestIndex = getBestIndex(specificStands)
		// Car `c` chooses stand with `bestIndex`
		c.waitForStandStarted = time.Now()
		specificStands[bestIndex].queue <- c
		// Car `c` has entered the back of the stand queue
	}
}

// Used as goroutine for every stand; simulates refueling and after that sends cars to cash registers
func (sor *standOrRegister) standFillAndDistributeToRegisters(registers []*standOrRegister) {
	var bestIndex int
	var c car
	for {
		c = <-sor.queue
		sor.isUsed = true
		// Car `c` is at the front of the stand queue => It starts refueling
		randSleepExcl(sor.minTime, sor.maxTime)
		bestIndex = getBestIndex(registers)
		// Car `c` chooses cash register with `bestIndex`
		c.waitForRegisterStarted = time.Now()
		registers[bestIndex].queue <- c
		// Car `c` has entered the back of the cash register queue
		sor.isUsed = false
	}
}

// Used as goroutine for every cash register; simulates paying and then makes cars leave the gas station
func (sor *standOrRegister) registerCashoutAndLeave() {
	var c car
	for {
		c = <-sor.queue
		sor.isUsed = true
		// Car `c` is at the front of the cash register queue => It starts paying
		randSleepExcl(sor.minTime, sor.maxTime)
		// Car `c` is leaving the gas station
		c.departureStarted = time.Now()
		go updateStats(c)
		sor.isUsed = false
	}
}

func main() {
	fmt.Println("Hello gas station!\nSimulation started ...")
	totalSimulationTime := time.Now()

	// Read data from config.yaml to the conf variable
	conf, err := readConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simulation
	// Make cash registers and start their goroutines that accept cars from stands and makes them leave the gas station after their payment is done
	registerConf := conf.Registers
	nRegisters := registerConf.Count
	registers := make([]*standOrRegister, nRegisters)
	for i := 0; i < nRegisters; i++ {
		registers[i] = &standOrRegister{i, false, make(chan car, registerConf.QueueLengthMax), registerConf.HandleTimeMin, registerConf.HandleTimeMax}
		go registers[i].registerCashoutAndLeave()
	}

	// Make allStands and start their goroutines that accept cars from sharedQueue, refuels them, and sends them to cash registers
	var nSpecificStands int

	fuelTypeConfs := map[int]standConfigRepresentation{
		gas:      standConfigRepresentation(conf.Stations.Gas),
		diesel:   standConfigRepresentation(conf.Stations.Diesel),
		lpg:      standConfigRepresentation(conf.Stations.LPG),
		electric: standConfigRepresentation(conf.Stations.Electric),
	}
	allStands := make(map[int][]*standOrRegister)
	for key, value := range fuelTypeConfs {
		nSpecificStands = value.Count
		allStands[key] = make([]*standOrRegister, nSpecificStands)
		for i := 0; i < nSpecificStands; i++ {
			allStands[key][i] = &standOrRegister{(key+1)*10 + i, false, make(chan car, value.QueueLengthMax), value.ServeTimeMin, value.ServeTimeMax}
			go allStands[key][i].standFillAndDistributeToRegisters(registers)
		}
	}

	// Make a sharedQueue and start a goroutine that sends cars to specific queues according to their fuel type
	carsConf := conf.Cars
	nCars := carsConf.Count
	sq := sharedQueue{make(chan car, carsConf.SharedQueueLengthMax)}
	go sq.distributeToStands(allStands)

	// Start spawning cars
	globalWG.Add(nCars)
	go carSpawner(nCars, carsConf.ArrivalTimeMin, carsConf.ArrivalTimeMax, &sq)
	globalWG.Wait() // Wait for all the cars to leave (also wait for stats update)

	fmt.Println("Simulation ended, took", time.Since(totalSimulationTime))
	finalizeStats() // set maximum values and compute average values

	// Print stats to terminal
	toPrint := fmt.Sprintf("%+v\n", globalStats)
	toPrint = strings.Replace(toPrint, " ", "\n", -1)
	toPrint = strings.Replace(toPrint, ":{", ":{\n", -1)
	toPrint = strings.Replace(toPrint, "}", "}\n", -1)
	fmt.Println(toPrint)

	// Write stats to yaml
	err = writeGlobalStats()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Stats saved to 'output.yaml'\nBye!")
}
