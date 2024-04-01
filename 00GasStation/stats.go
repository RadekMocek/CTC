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
	timeSpentWaitingForAndInSharedQueue := c.waitForStandStarted.Sub(c.waitForSharedQueueStarted)
	timeSpenWaitingForAndRefueling := c.waitForRegisterStarted.Sub(c.waitForStandStarted)
	timeSpentWaitingForAndPaying := c.departureStarted.Sub(c.waitForRegisterStarted)

	globalStatsLock.Lock()
	// Update Totals
	// - SharedQueue
	globalStats.SharedQueue.TotalCars++
	globalStats.SharedQueue.TotalTime += timeSpentWaitingForAndInSharedQueue
	// - SpecificQueue
	switch c.fuelType {
	case gas:
		globalStats.Stations.Gas.TotalCars++
		globalStats.Stations.Gas.TotalTime += timeSpenWaitingForAndRefueling
	case diesel:
		globalStats.Stations.Diesel.TotalCars++
		globalStats.Stations.Diesel.TotalTime += timeSpenWaitingForAndRefueling
	case lpg:
		globalStats.Stations.LPG.TotalCars++
		globalStats.Stations.LPG.TotalTime += timeSpenWaitingForAndRefueling
	case electric:
		globalStats.Stations.Electric.TotalCars++
		globalStats.Stations.Electric.TotalTime += timeSpenWaitingForAndRefueling
	}
	// - RegisterQueue
	globalStats.Registers.TotalCars++
	globalStats.Registers.TotalTime += timeSpentWaitingForAndPaying
	// Update Maximums
	// - SharedQueue
	if timeSpentWaitingForAndInSharedQueue > globalMaxShared {
		globalMaxShared = timeSpentWaitingForAndInSharedQueue
	}
	// - SpecificQueue
	if timeSpenWaitingForAndRefueling > globalMaxSpecific[c.fuelType] {
		globalMaxSpecific[c.fuelType] = timeSpenWaitingForAndRefueling
	}
	// - RegisterQueue
	if timeSpentWaitingForAndPaying > globalMaxRegisters {
		globalMaxRegisters = timeSpentWaitingForAndPaying
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
