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
	file, _ := os.Open("./paths.csv")
	defer file.Close()

	ch := make(chan map[string]float64, 10)
	wg := sync.WaitGroup{}
	for chunk := range ingress.Read(file, 10*10*1024) {
		wg.Add(1)
		go func(chunk []byte) {
			estimates := fare.Aggregate(chunk)
			ch <- estimates
			wg.Done()
		}(chunk)
	}

	wg2 := sync.WaitGroup{}
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	wg2.Add(1)
	go func() {
		fmt.Println("Setting up writer")
		for estimates := range ch {
			for id, estimate := range estimates {
				err = writer.Write([]string{id, fmt.Sprintf("%.02f", estimate)})
				if err != nil {
					fmt.Println("err", err)
				}
			}
		}
		wg2.Done()
	}()

	wg.Wait()
	close(ch)
	wg2.Wait()
}
