package fare

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kasperbe/beat/gps"
)

// calculatePrice for a distance travelled.
func calculatePrice(a, b gps.Coordinate, time float64) (float64, error) {

	distance := gps.Distance(a, b)
	mps := distance / float64(time)
	speed := mps * 3.6
	km := distance / 1000

	if speed > 100 {
		return 0, errors.New("Moving too fast")
	}

	if speed > 10 {
		return km * 1.30, nil
	}

	if speed <= 10 {
		return float64(time/60/60) * 11.90, nil
	}

	return 0, nil
}

// SegmentPoint holds the last known position of a segment
type SegmentPoint struct {
	ID         string
	Coordinate gps.Coordinate
	Timestamp  float64
}

// New Aggregator
func New() Aggregator {
	return Aggregator{
		fares:     make(map[string]*SegmentPoint),
		Estimates: make(map[string]float64),
	}
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
func (aggr *Aggregator) AddSegment(id string, coord gps.Coordinate, timestamp float64) {
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

	time := timestamp - fare.Timestamp
	price, err := calculatePrice(fare.Coordinate, coord, time)
	if err != nil {
		// Price is above 100 km/h
		// Segment is not valid
		return
	}

	fare.Timestamp = 0
	fare.Coordinate.Lat = 0
	fare.Coordinate.Lng = 0

	aggr.Estimates[id] += price
}

// Aggregate segments to the estimated cost
func Aggregate(chunk []byte) map[string]float64 {
	aggr := New()

	lines := strings.Split(string(chunk), "\n")
	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		row := strings.Split(line, ",")
		id := row[0]
		lat, _ := strconv.ParseFloat(row[1], 64)
		lng, _ := strconv.ParseFloat(row[2], 64)
		timestamp, _ := strconv.ParseFloat(row[3], 64)
		aggr.AddSegment(id, gps.Coordinate{Lat: lat, Lng: lng}, timestamp)
	}

	return aggr.Estimates
}
