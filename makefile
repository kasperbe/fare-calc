profile:
	go build . && ./beat -profile=profile.prof -input=./paths-big.csv
	go tool pprof profile.prof

run:
	go build . && ./beat