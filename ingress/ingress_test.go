package ingress

import (
	"testing"
)

func TestIndexLastSegment(t *testing.T) {
	tt := []struct {
		Name     string
		Input    []byte
		Expected int
	}{
		{"No newline", []byte(`1,some,thing,else`), 17},
		{"Same ID", []byte(string("1,some,thing,else\n1,some,thing,else\n1,some,thing,else\n1,some,thing,else")), 53},
		{"Diff ids", []byte(string("1,some,thing,else\n2,some,thing,else\n3,some,thing,else")), 35},
		{"Empty buffer", []byte{}, 0},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			offset := indexLastSegment(tc.Input)
			if offset != tc.Expected {
				t.Errorf("Wrong offset: expected %d, got %d", tc.Expected, offset)
			}
		})
	}
}

func TestRead(t *testing.T) {
	input := []byte(`1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
2,some,thing,else
2,some,thing,else`)

	offset := indexLastSegment(input)
	if offset != 107 {
		t.Error("Wrong offset")
	}
}
