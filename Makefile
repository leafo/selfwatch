.PHONY: build install test

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS := -X 'main.commitHash=$(COMMIT_HASH)' -X 'main.buildDate=$(BUILD_DATE)'

build:
	go build -ldflags "$(LDFLAGS)" -o bin/selfwatch

install:
	go install -ldflags "$(LDFLAGS)"

test:
	go test -v ./...
