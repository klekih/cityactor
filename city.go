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

// Junction is the message send back from the city with information
// about a junction
type Junction struct {
	Loc Location
}

const (
	// SendReport is a message passed from an actor to the city
	// indicating its status (e.g. location).
	SendReport = iota

	// AskForJunction is a message passed from an actor to the city.
	// A response is awaited.
	AskForJunction = iota

	// RespondWithJunction is a message passed from the city to
	// an actor and it contains junction data.
	RespondWithJunction = iota
)

// Envelope is the container for different messages sent back
// and forth between an actor and a city
type Envelope struct {
	MessageType int
	Payload     interface{}
}

// Connect is the typical method used for connecting to
// a city.
func Connect() (chan Report, chan Junction) {

	var sendReportChan = make(chan Report)
	var junctionChan = make(chan Junction)

	go func() {
		for {
			select {
			case r := <-sendReportChan:
				conn, err := net.Dial("tcp", "localhost:7450")
				if err != nil {
					fmt.Println("Error on dialing", err)
					break
				}
				defer conn.Close()

				env := Envelope{
					MessageType: SendReport,
					Payload:     r}

				gob.Register(r)
				enc := gob.NewEncoder(conn)
				err = enc.Encode(env)
				if err != nil {
					fmt.Println("Error on sending data", err)
				}
			case j := <-junctionChan:
				conn, err := net.Dial("tcp", "localhost:7450")
				if err != nil {
					fmt.Println("Error on dialing", err)
					break
				}
				defer conn.Close()
				env := Envelope{
					MessageType: AskForJunction,
					Payload:     j}

				gob.Register(Junction{})
				enc := gob.NewEncoder(conn)
				err = enc.Encode(env)
				if err != nil {
					fmt.Println("Error on sending data", err)
				}

				dec := gob.NewDecoder(conn)
				env = Envelope{}
				err = dec.Decode(&env)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Received response with junction", env)
				junctionChan <- Junction{}
			}
		}
	}()

	return sendReportChan, junctionChan
}
