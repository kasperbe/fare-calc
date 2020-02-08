package fare

import (
	"math"
	"testing"

	"github.com/kasperbe/beat/gps"
)

func TestPriceCalculation(t *testing.T) {
	tt := []struct {
		Name          string
		PointA        gps.Coordinate
		PointB        gps.Coordinate
		Starttime     int64
		Endtime       int64
		Expected      float64
		ExpectedError bool
	}{
		{"Moving 1 hour day", gps.Coordinate{Lat: 55.403755, Lng: 10.402370}, gps.Coordinate{Lat: 56.158150, Lng: 10.212030}, 1405591384, 1405594984, 62.77, false},
		{"Moving 1 hour night", gps.Coordinate{Lat: 55.403755, Lng: 10.402370}, gps.Coordinate{Lat: 56.158150, Lng: 10.212030}, 1581120000, 1581123600, 110.27, false},
		{"Idle 1 hour seconds", gps.Coordinate{Lat: 55.403755, Lng: 10.402370}, gps.Coordinate{Lat: 55.403755, Lng: 10.402370}, 1405591384, 1405594984, 11.90, false},
		{"Path example", gps.Coordinate{Lat: 37.966613, Lng: 23.728375}, gps.Coordinate{Lat: 37.966203, Lng: 23.728597}, 1405594984, 1405594992, 0.04, false},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := calculatePrice(tc.PointA, tc.PointB, tc.Starttime, tc.Endtime)
			if err != nil {
				if tc.ExpectedError != true {
					t.Errorf("Got speed above 100km/h, but did not expect it to be.")
				}
			}

			if math.Round(result*100)/100 != tc.Expected {
				t.Errorf("Expected %0.2f got %0.2f", tc.Expected, result)
			}
		})
	}
}

func TestSegmentAggregation(t *testing.T) {
	tt := []struct {
		Name     string
		Chunk    []byte
		Expected map[string]float64
	}{
		{"Minimum ride cost", []byte(string("1,55.403755,10.402370,2\n1,55.403755,10.402370,2")), map[string]float64{"1": 1.30}},
		{"Idle ride", []byte(string("1,55.403755,10.402370,1581161779\n1,55.403755,10.402370,1581165373")), map[string]float64{"1": 13.18}},
		{"Multi idle ride", []byte(string("1,55.403755,10.402370,1581161779\n1,55.403755,10.402370,1581165373\n2,55.403755,10.402370,1581161779\n2,55.403755,10.402370,1581165373")), map[string]float64{"1": 13.18, "2": 13.18}},
		{"empty chunk", []byte{}, map[string]float64{}},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			estimates := Aggregate(tc.Chunk)
			for id, estimate := range estimates {
				result := math.Round(estimate*100) / 100

				if expected, ok := tc.Expected[id]; !ok {
					t.Errorf("Got estimate for id: %s but was not expected", id)
				} else if result != expected {
					t.Errorf("Got id: %s, expected estimate: %.02f got %.02f", id, expected, result)
				}
			}
		})
	}
}
