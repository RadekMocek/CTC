package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello gas station!")

	conf, err := readConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(conf.Stations.Gas.ServeTimeMax)
}
