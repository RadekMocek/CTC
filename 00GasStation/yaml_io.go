package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type gasStationConfig struct {
	Cars struct {
		Count                int `yaml:"count"`
		ArrivalTimeMin       int `yaml:"arrival_time_min"`
		ArrivalTimeMax       int `yaml:"arrival_time_max"`
		SharedQueueLengthMax int `yaml:"shared_queue_length_max"`
	}
	Stations struct {
		Gas struct {
			Count          int `yaml:"count"`
			ServeTimeMin   int `yaml:"serve_time_min"`
			ServeTimeMax   int `yaml:"serve_time_max"`
			QueueLengthMax int `yaml:"queue_length_max"`
		}
		Diesel struct {
			Count          int `yaml:"count"`
			ServeTimeMin   int `yaml:"serve_time_min"`
			ServeTimeMax   int `yaml:"serve_time_max"`
			QueueLengthMax int `yaml:"queue_length_max"`
		}
		LPG struct {
			Count          int `yaml:"count"`
			ServeTimeMin   int `yaml:"serve_time_min"`
			ServeTimeMax   int `yaml:"serve_time_max"`
			QueueLengthMax int `yaml:"queue_length_max"`
		}
		Electric struct {
			Count          int `yaml:"count"`
			ServeTimeMin   int `yaml:"serve_time_min"`
			ServeTimeMax   int `yaml:"serve_time_max"`
			QueueLengthMax int `yaml:"queue_length_max"`
		}
	}
	Registers struct {
		Count          int `yaml:"count"`
		HandleTimeMin  int `yaml:"handle_time_min"`
		HandleTimeMax  int `yaml:"handle_time_max"`
		QueueLengthMax int `yaml:"queue_length_max"`
	}
}

// https://github.com/go-yaml/yaml?tab=readme-ov-file#example
func readConfig(filename string) (*gasStationConfig, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &gasStationConfig{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return conf, err
}

type gasStationStats struct {
	SharedQueue struct {
		TotalCars int           `yaml:"total_cars"`
		TotalTime time.Duration `yaml:"total_time"`
		AvgTime   time.Duration `yaml:"avg_time"`
		MaxTime   time.Duration `yaml:"max_time"`
	}
	Stations struct {
		Gas struct {
			TotalCars int           `yaml:"total_cars"`
			TotalTime time.Duration `yaml:"total_time"`
			AvgTime   time.Duration `yaml:"avg_time"`
			MaxTime   time.Duration `yaml:"max_time"`
		}
		Diesel struct {
			TotalCars int           `yaml:"total_cars"`
			TotalTime time.Duration `yaml:"total_time"`
			AvgTime   time.Duration `yaml:"avg_time"`
			MaxTime   time.Duration `yaml:"max_time"`
		}
		LPG struct {
			TotalCars int           `yaml:"total_cars"`
			TotalTime time.Duration `yaml:"total_time"`
			AvgTime   time.Duration `yaml:"avg_time"`
			MaxTime   time.Duration `yaml:"max_time"`
		}
		Electric struct {
			TotalCars int           `yaml:"total_cars"`
			TotalTime time.Duration `yaml:"total_time"`
			AvgTime   time.Duration `yaml:"avg_time"`
			MaxTime   time.Duration `yaml:"max_time"`
		}
	}
	Registers struct {
		TotalCars int           `yaml:"total_cars"`
		TotalTime time.Duration `yaml:"total_time"`
		AvgTime   time.Duration `yaml:"avg_time"`
		MaxTime   time.Duration `yaml:"max_time"`
	}
}

// https://stackoverflow.com/a/65207714
func writeGlobalStats() (returnValue error) {
	file, err := os.OpenFile("output.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening/creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			returnValue = fmt.Errorf("error closing file: %v", err)
		}
	}(file)
	enc := yaml.NewEncoder(file)
	err = enc.Encode(globalStats)
	if err != nil {
		return fmt.Errorf("error encoding: %v", err)
	}
	return nil
}
