profile:
	go build . && ./beat
	go tool pprof x.pproff

run:
	go build . && ./beat