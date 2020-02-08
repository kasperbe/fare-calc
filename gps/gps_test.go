package gps

import (
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tt := []struct {
		Name     string
		A        Coordinate
		B        Coordinate
		Expected float64
	}{
		{"Aarhus-Odense", Coordinate{55.403755, 10.402370}, Coordinate{56.158150, 10.212030}, 84819.406710},
		{"Odense-Paris", Coordinate{55.396229, 10.390600}, Coordinate{48.856613, 2.352222}, 910759.517940},
	}

	for _, tc := range tt {
		result := math.Round(Distance(tc.A, tc.B)*100000) / 100000

		if result != tc.Expected {
			t.Errorf("Error in [%s] calculating distance. Expected %f, got %f", tc.Name, tc.Expected, result)
		}
	}
}
