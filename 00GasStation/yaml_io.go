package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type gasStationConfig struct {
	Cars struct {
		Count          int `yaml:"count"`
		ArrivalTimeMin int `yaml:"arrival_time_min"`
		ArrivalTimeMax int `yaml:"arrival_time_max"`
	}
	Stations struct {
		Gas struct {
			Count        int `yaml:"count"`
			ServeTimeMin int `yaml:"serve_time_min"`
			ServeTimeMax int `yaml:"serve_time_max"`
		}
		Diesel struct {
			Count        int `yaml:"count"`
			ServeTimeMin int `yaml:"serve_time_min"`
			ServeTimeMax int `yaml:"serve_time_max"`
		}
		LPG struct {
			Count        int `yaml:"count"`
			ServeTimeMin int `yaml:"serve_time_min"`
			ServeTimeMax int `yaml:"serve_time_max"`
		}
		Electric struct {
			Count        int `yaml:"count"`
			ServeTimeMin int `yaml:"serve_time_min"`
			ServeTimeMax int `yaml:"serve_time_max"`
		}
	}
	Registers struct {
		Count          int `yaml:"count"`
		ArrivalTimeMin int `yaml:"handle_time_min"`
		ArrivalTimeMax int `yaml:"handle_time_max"`
	}
}

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
