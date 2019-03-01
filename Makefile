.PHONY: bin
bin: main.go ubcrmp/*.go
	rm -rf data/*
	go build -o bin/ubc-rmp-data main.go

