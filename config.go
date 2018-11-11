package main

import "encoding/json"
import "os"
import "fmt"

type Config struct {
	SimulationStep int `json:"simulationStep"` // given in milliseconds
}

func getConfig() *Config {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error in reading configuration file ", err)
	}

	return &config
}
