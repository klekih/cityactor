package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func generateRandomRoute(config *Config) (rt *Route, err error) {

	link := "http://localhost:8989/route?point=46.748654%2C23.535461&point=46.792942%2C23.664862&locale=en-US&vehicle=car&weighting=fastest&elevation=false&use_miles=false&layer=Omniscale&points_encoded=false"

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
