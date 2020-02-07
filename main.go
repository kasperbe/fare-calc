package main

import (
	"fmt"
	"os"

	"github.com/kasperbe/beat/ingress"
)

func main() {
	file, _ := os.Open("../beat-test/paths-big.csv")

	a := 0
	for chunk := range ingress.Read(file, 10*10*1024) {
		a += len(chunk)
		//fmt.Println(string(chunk))
	}

	fmt.Println(a)
}
