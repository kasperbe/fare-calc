package ingress

import (
	"fmt"
	"os"
	"testing"
)

func TestFindDelimiterWithoutNewline(t *testing.T) {
	input := []byte(`1,some,thing,else`)

	offset := findDelimiter(input)
	if offset != 17 {
		t.Errorf("Wrong offset: expected %d, got %d", 17, offset)
	}
}

func TestFindDelimiter(t *testing.T) {
	input := []byte(`1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
1,some,thing,else
2,some,thing,else
2,some,thing,else`)

	offset := findDelimiter(input)
	if offset != 107 {
		t.Error("Wrong offset")
	}
}
func TestReadFile(t *testing.T) {
	file, _ := os.Open("./paths.csv")

	for chunk := range Read(file, 10*2048) {
		fmt.Println("Recieving: ", string(chunk))
	}

	t.Log()
	t.Error()

}
