.PHONY: install test

install:
	go install github.com/leafo/selfwatch

test:
	go test -v github.com/leafo/selfwatch/selfwatch
