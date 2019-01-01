package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Route is the json equivalent of the graphhopper answer for
// computing routes
type Route struct {
	Hints struct {
		VisitedNodesAverage string `json:"visited_nodes.average"`
		VisitedNodesSum     string `json:"visited_nodes.sum"`
	} `json:"hints"`
	Info struct {
		Copyrights []string `json:"copyrights"`
		Took       int      `json:"took"`
	} `json:"info"`
	Paths []struct {
		Distance      float64   `json:"distance"`
		Weight        float64   `json:"weight"`
		Time          int       `json:"time"`
		Transfers     int       `json:"transfers"`
		PointsEncoded bool      `json:"points_encoded"`
		Bbox          []float64 `json:"bbox"`
		Points        struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"points"`
		Instructions []struct {
			Distance    float64 `json:"distance"`
			Heading     float64 `json:"heading,omitempty"`
			Sign        int     `json:"sign"`
			Interval    []int   `json:"interval"`
			Text        string  `json:"text"`
			Time        int     `json:"time"`
			StreetName  string  `json:"street_name"`
			ExitNumber  int     `json:"exit_number,omitempty"`
			Exited      bool    `json:"exited,omitempty"`
			TurnAngle   float64 `json:"turn_angle,omitempty"`
			LastHeading float64 `json:"last_heading,omitempty"`
		} `json:"instructions"`
		Legs    []interface{} `json:"legs"`
		Details struct {
		} `json:"details"`
		Ascend           int `json:"ascend"`
		Descend          int `json:"descend"`
		SnappedWaypoints struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"snapped_waypoints"`
	} `json:"paths"`
}

// Point is a simple structure for a map point with longitude and latitude
type Point struct {
	long float64
	lat  float64
}

func generateRandomRoute(config *Config) (rt *Route, err error) {

	// Choose random points from the bounding box.
	// The bounding box should have two points which represent
	// the opposite corners of a square. Two random points from within
	// this square are choosed to define a route's start and finish.

	// default values if something goes wrong
	p1CoordLong := 46.748654
	p1CoordLat := 23.53546
	p2CoordLong := 46.792942
	p2CoordLat := 23.664862

	first, err := parseCoordinates(config.BoundingBox.First)
	second, err := parseCoordinates(config.BoundingBox.Second)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	p1CoordLong = first.long + r.Float64()*(second.long-first.long)
	p1CoordLat = first.lat + r.Float64()*(second.lat-first.lat)

	p2CoordLong = first.long + r.Float64()*(second.long-first.long)
	p2CoordLat = first.lat + r.Float64()*(second.lat-first.lat)

	link := fmt.Sprintf("%s%f%s%f%s%f%s%f%s",
		"http://localhost:8989/route?point=",
		p1CoordLong,
		"%2C",
		p1CoordLat,
		"&point=",
		p2CoordLong,
		"%2C",
		p2CoordLat,
		"&locale=en-US&vehicle=car&weighting=fastest&elevation=false&use_miles=false&layer=Omniscale&points_encoded=false")

	fmt.Println("Tyring link:", link)

	retries := config.GetRouteRetries

	resp, err := http.Get(link)

	for (retries > 0) && (err != nil) {
		fmt.Println("No valid link for getting route, retrying", err)

		resp, err = http.Get(link)

		retries--
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var route Route
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &route)

	return &route, nil
}

func parseCoordinates(strCoord string) (pt *Point, err error) {
	arrCoord := strings.Split(strCoord, ",")
	if len(arrCoord) < 2 {
		fmt.Println("Bad coordinates on", arrCoord)
		return nil, errors.New("Too few arguments in")
	}
	long, err1 := strconv.ParseFloat(arrCoord[0], 64)
	lat, err2 := strconv.ParseFloat(arrCoord[1], 64)

	var finalErrStr string

	if err1 != nil {
		finalErrStr = fmt.Sprintf("%s%s", finalErrStr, err1)
	}

	if err2 != nil {
		finalErrStr = fmt.Sprintf("%s%s", finalErrStr, err2)
	}

	var finalErr error
	if finalErrStr != "" {
		finalErr = errors.New(finalErrStr)
	}
	return &Point{long: long, lat: lat}, finalErr
}
