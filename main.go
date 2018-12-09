package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, world. I'm a city actor.")

	// load the config file
	config := getConfig()
	fmt.Println("Using config: ", config)

	// now start the "walker"
	// (the timer which advances our position)
	walkerChan := StartWalker(config)

	finished := false
	for finished != true {
		status := <-walkerChan
		finished = status.arrivedAtDestination
		fmt.Println("walker done: ", finished)
	}

	fmt.Println("I'm out of job")
	os.Exit(0)
}
