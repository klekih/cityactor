package main

import (
	"fmt"
	"time"
)

var ticker *time.Ticker
var myRoute *Route

// WalkerStatus is the answer for reporting walker status
type WalkerStatus struct {
	arrivedAtDestination bool
}

// StartWalker is the main entry point for starting a walker
func StartWalker(period int) chan WalkerStatus {

	// first, get a route to go on
	myRoute = getRoute()

	// create the chan to report back the status
	c := make(chan WalkerStatus)

	duration := time.Duration(period)
	ticker = time.NewTicker(duration * time.Millisecond)

	go func() {
		for t := range ticker.C {
			fmt.Println("tick at ", t)
			advance()
		}
	}()

	return c
}

func getRoute() *Route {

	route, err := generateRandomRoute()

	if err != nil {
		panic(err)
	}

	return route
}

func advance() {

	if myRoute == nil {
		panic("No route to go on")
	}

	if len(myRoute.Paths) == 0 {
		panic("Route has no path to go on")
	}

}
