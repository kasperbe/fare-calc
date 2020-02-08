package fare

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kasperbe/beat/gps"
)

const (
	day   = 0.75
	night = 1.30
	idle  = 11.90
)

// calculatePrice for a distance travelled.
// We assume that if the segment is started at either day or night
// it also ends in the same time range.
func calculatePrice(a, b gps.Coordinate, starttime, endtime int64) (float64, error) {
	duration := float64(endtime - starttime)

	h, _, _ := time.Unix(starttime, 0).Clock()
	distance := gps.Distance(a, b)
	mps := distance / duration
	speed := mps * 3.6
	km := distance / 1000

	if speed > 100 {
		return 0, errors.New("Moving too fast")
	}

	if speed > 10 {
		if h >= 5 && h <= 24 {
			return km * 0.74, nil
		}
		return km * 1.30, nil
	}

	if speed <= 10 {
		return duration / 60 / 60 * 11.90, nil
	}

	return 0, nil
}

// SegmentPoint holds the last known position of a segment
type SegmentPoint struct {
	ID         string
	Coordinate gps.Coordinate
	Timestamp  int64
}

// Aggregator handles
type Aggregator struct {
	Estimates map[string]float64
	fares     map[string]*SegmentPoint
}

// AddSegment checks if an existing point already exists.
// If an existing point exists the points are aggregated and the estimated segment cost is
// added to the fare total.
//
// After a pair of points has been matched, the last known position is reset to avoid double counting
// since two consecutive tuples belonging to the same ride form a segment.
//
// If the speed is above 100 km/h an error occurs in the calculate price function
// we simply skip the point and move forward.
func (aggr *Aggregator) AddSegment(id string, coord gps.Coordinate, timestamp int64) {
	fare, exists := aggr.fares[id]
	if !exists {
		aggr.fares[id] = &SegmentPoint{
			ID:         id,
			Coordinate: coord,
			Timestamp:  timestamp,
		}

		aggr.Estimates[id] = 1.30

		return
	}

	if fare.Coordinate.Lat == 0 {
		fare.Coordinate = coord
		fare.Timestamp = timestamp

		return
	}

	price, err := calculatePrice(fare.Coordinate, coord, fare.Timestamp, timestamp)
	if err != nil {
		// Price is above 100 km/h
		// Segment is not valid
		return
	}

	// Reuse the allocated SegmentPoint for the id
	fare.Timestamp = 0
	fare.Coordinate.Lat = 0
	fare.Coordinate.Lng = 0

	aggr.Estimates[id] += price
}

// Aggregate segments to the estimated cost
func Aggregate(chunk []byte) map[string]float64 {
	if len(chunk) < 1 {
		return map[string]float64{}
	}

	aggr := Aggregator{
		fares:     make(map[string]*SegmentPoint),
		Estimates: make(map[string]float64),
	}

	lines := strings.Split(string(chunk), "\n")
	var lat, lng float64
	var timestamp int64
	for _, line := range lines {
		if len(line) < 1 {
			continue
		}
		row := strings.Split(line, ",")
		lat, _ = strconv.ParseFloat(row[1], 64)
		lng, _ = strconv.ParseFloat(row[2], 64)
		timestamp, _ = strconv.ParseInt(row[3], 10, 64)
		aggr.AddSegment(row[0], gps.Coordinate{Lat: lat, Lng: lng}, timestamp)
	}

	return aggr.Estimates
}
