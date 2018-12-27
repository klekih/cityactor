package main

import "encoding/json"
import "os"
import "fmt"

// Config is the corresponding go structure of the json configuration file
type Config struct {
	SimulationStep  int `json:"simulationStep"`
	GetRouteRetries int `json:"getRouteRetries"`
	BoundingBox     struct {
		First  string `json:"first"`
		Second string `json:"second"`
		Third  string `json:"third"`
		Fourth string `json:"fourth"`
	} `json:"boundingBox"`
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
