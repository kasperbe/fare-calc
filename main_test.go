package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestIntegration(t *testing.T) {
	outname := "test.csv"
	outfilename = &outname

	results := map[string]float64{
		"8": 5.32,
		"1": 7.52,
		"2": 7.13,
		"3": 25.08,
		"5": 12.85,
		"6": 5.00,
		"4": 2.00,
		"7": 20.47,
		"9": 3.91,
	}

	main()

	f, _ := os.Open(*outfilename)
	data, _ := csv.NewReader(f).ReadAll()

	if len(results) != len(data) {
		t.Errorf("Got too many lines in result. Expected: %d, got %d", len(results), len(data))
	}
	for _, row := range data {
		estimate, _ := strconv.ParseFloat(row[1], 64)
		if results[row[0]] != estimate {
			t.Errorf("Wrong price provided. Expected %.02f, got %.02f", results[row[0]], estimate)
		}
	}

	os.Remove(*outfilename)
}
