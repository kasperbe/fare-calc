profile:
	go build . && ./beat -profile=profile.prof -input=./paths.csv
	go tool pprof profile.prof

run:
	go build . && ./beat