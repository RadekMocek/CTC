package main

import (
	"sync"
	"time"
)

var globalStatsLock sync.Mutex
var globalStats gasStationStats

var globalMaxShared, globalMaxRegisters time.Duration = 0, 0
var globalMaxSpecific = map[int]time.Duration{gas: 0, diesel: 0, lpg: 0, electric: 0}

func updateStats(c *car) {
	// Calculate durations
	timeSpentInSharedQueue := c.standQueueEnteredTime.Sub(c.sharedQueueEnteredTime)
	if timeSpentInSharedQueue < 0 { // Because of time imprecisions
		timeSpentInSharedQueue = 0
	}
	timeSpentInSpecificQueueAndRefueling := c.refuelingFinishedTime.Sub(c.standQueueEnteredTime)
	if timeSpentInSpecificQueueAndRefueling < 0 {
		timeSpentInSpecificQueueAndRefueling = 0
	}
	timeSpentInRegisterQueueAndPaying := c.exitTime.Sub(c.refuelingFinishedTime)
	//fmt.Println(timeSpentInRegisterQueueAndPaying, c.exitTime, c.refuelingFinishedTime)
	if timeSpentInRegisterQueueAndPaying < 0 {
		timeSpentInRegisterQueueAndPaying = 0
	}
	//fmt.Println(timeSpentInRegisterQueueAndPaying)
	//fmt.Println()

	globalStatsLock.Lock()
	// Update Totals
	// - SharedQueue
	globalStats.SharedQueue.TotalCars++
	globalStats.SharedQueue.TotalTime += timeSpentInSharedQueue
	// - SpecificQueue
	switch c.fuelType {
	case gas:
		globalStats.Stations.Gas.TotalCars++
		globalStats.Stations.Gas.TotalTime += timeSpentInSpecificQueueAndRefueling
	case diesel:
		globalStats.Stations.Diesel.TotalCars++
		globalStats.Stations.Diesel.TotalTime += timeSpentInSpecificQueueAndRefueling
	case lpg:
		globalStats.Stations.LPG.TotalCars++
		globalStats.Stations.LPG.TotalTime += timeSpentInSpecificQueueAndRefueling
	case electric:
		globalStats.Stations.Electric.TotalCars++
		globalStats.Stations.Electric.TotalTime += timeSpentInSpecificQueueAndRefueling
	}
	// - RegisterQueue
	globalStats.Registers.TotalCars++
	globalStats.Registers.TotalTime += timeSpentInRegisterQueueAndPaying
	// Update Maximums
	// - SharedQueue
	if timeSpentInSharedQueue > globalMaxShared {
		globalMaxShared = timeSpentInSharedQueue
	}
	// - SpecificQueue
	if timeSpentInSpecificQueueAndRefueling > globalMaxSpecific[c.fuelType] {
		globalMaxSpecific[c.fuelType] = timeSpentInSpecificQueueAndRefueling
	}
	// - RegisterQueue
	if timeSpentInRegisterQueueAndPaying > globalMaxRegisters {
		globalMaxRegisters = timeSpentInRegisterQueueAndPaying
	}
	globalStatsLock.Unlock()

	globalWG.Done() // Car definitely left the gas station
}

func finalizeStats() {
	globalStats.SharedQueue.MaxTime = globalMaxShared
	globalStats.Stations.Gas.MaxTime = globalMaxSpecific[gas]
	globalStats.Stations.Diesel.MaxTime = globalMaxSpecific[diesel]
	globalStats.Stations.LPG.MaxTime = globalMaxSpecific[lpg]
	globalStats.Stations.Electric.MaxTime = globalMaxSpecific[electric]
	globalStats.Registers.MaxTime = globalMaxRegisters

	if globalStats.SharedQueue.TotalCars > 0 {
		globalStats.SharedQueue.AvgTime = globalStats.SharedQueue.TotalTime / time.Duration(globalStats.SharedQueue.TotalCars)
	}
	if globalStats.Stations.Gas.TotalCars > 0 {
		globalStats.Stations.Gas.AvgTime = globalStats.Stations.Gas.TotalTime / time.Duration(globalStats.Stations.Gas.TotalCars)
	}
	if globalStats.Stations.Diesel.TotalCars > 0 {
		globalStats.Stations.Diesel.AvgTime = globalStats.Stations.Diesel.TotalTime / time.Duration(globalStats.Stations.Diesel.TotalCars)
	}
	if globalStats.Stations.LPG.TotalCars > 0 {
		globalStats.Stations.LPG.AvgTime = globalStats.Stations.LPG.TotalTime / time.Duration(globalStats.Stations.LPG.TotalCars)
	}
	if globalStats.Stations.Electric.TotalCars > 0 {
		globalStats.Stations.Electric.AvgTime = globalStats.Stations.Electric.TotalTime / time.Duration(globalStats.Stations.Electric.TotalCars)
	}
	if globalStats.Registers.TotalCars > 0 {
		globalStats.Registers.AvgTime = globalStats.Registers.TotalTime / time.Duration(globalStats.Registers.TotalCars)
	}
}
