package main

import (
	"fmt"
	"sync"
	"time"
)

var mainWG sync.WaitGroup

func carSpawner(nCars int, minTime int, maxTime int) {
	for i := 0; i < nCars; i++ {
		time.Sleep(time.Duration(randRangeExcl(minTime, maxTime)) * time.Second)
		fmt.Println(i)
		mainWG.Done()
	}
}

func main() {
	fmt.Println("Hello gas station!")

	conf, err := readConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	nCars := conf.Cars.Count
	mainWG.Add(nCars)

	go carSpawner(nCars, conf.Cars.ArrivalTimeMin, conf.Cars.ArrivalTimeMax)

	mainWG.Wait()
	fmt.Println("Simulation done.")
}
