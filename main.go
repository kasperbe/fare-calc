package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/kasperbe/beat/fare"
	"github.com/kasperbe/beat/ingress"
)

var cpuprofile = flag.String("profile", "", "filename of profile output")
var inputfilename = flag.String("input", "./paths.csv", "filename of input csv file")
var outfilename = flag.String("output", "result.csv", "filename of result csv file")

func profile() func() {
	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		stop := profile()
		defer stop()
	}

	file, _ := os.Open(*inputfilename)

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
	file.Close()

	outfile, _ := os.Create(*outfilename)
	defer file.Close()
	writer := csv.NewWriter(outfile)
	defer writer.Flush()

	writewg := sync.WaitGroup{}
	defer writewg.Wait()

	writewg.Add(1)
	go func() {
		for estimates := range ch {
			for id, estimate := range estimates {
				err := writer.Write([]string{id, fmt.Sprintf("%.02f", estimate)})
				if err != nil {
					log.Fatalf("%w while trying to write file", err)
				}
			}
		}

		writewg.Done()
	}()

	readwg.Wait()
	close(ch)
}
