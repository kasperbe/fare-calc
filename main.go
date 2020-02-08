package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kasperbe/beat/fare"
	"github.com/kasperbe/beat/ingress"
)

func main() {
	file, _ := os.Open("./paths-big.csv")
	defer file.Close()

	ch := make(chan map[string]float64, 10)
	readwg := sync.WaitGroup{}
	for chunk := range ingress.Read(file, 10*10*1024) {
		readwg.Add(1)
		go func(chunk []byte) {
			estimates := fare.Aggregate(chunk)
			ch <- estimates
			readwg.Done()
		}(chunk)
	}

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writewg := sync.WaitGroup{}
	defer writewg.Wait()

	writewg.Add(1)
	go func() {
		for estimates := range ch {
			for id, estimate := range estimates {
				err = writer.Write([]string{id, fmt.Sprintf("%.02f", estimate)})
				if err != nil {
					fmt.Println("err", err)
				}
			}
		}
	}()

	readwg.Wait()
	close(ch)
}
