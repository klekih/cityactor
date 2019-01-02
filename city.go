package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

// CityInterface is a representation to a city entity
type CityInterface interface {
	SendVector()
	RetrieveJunction()
}

// Location aggregates the information about the place
// where an actor is at a certain moment in time
type Location struct {
	Long float64
	Lat  float64
}

// Report is the base type for reporting status and vectors
// to a city entity
type Report struct {
	Loc Location
}

// Connect is the typical method used for connecting to
// a city.
func Connect() chan Report {

	var sendChan = make(chan Report)

	go func() {
		for {
			select {
			case r := <-sendChan:
				conn, err := net.Dial("tcp", "localhost:7450")
				if err != nil {
					fmt.Println("Error on dialing", err)
					break
				}
				defer conn.Close()
				enc := gob.NewEncoder(conn)
				err = enc.Encode(r)
				if err != nil {
					fmt.Println("Error on sending data", err)
				}
			}
		}
	}()

	return sendChan
}
