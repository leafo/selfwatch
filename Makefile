.PHONY: build test

build:
	go build -o bin/selfwatch

test:
	go test -v ./...
