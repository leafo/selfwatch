.PHONY: build install test bundle release

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS := -X 'main.commitHash=$(COMMIT_HASH)' -X 'main.buildDate=$(BUILD_DATE)'

bundle:
	npm run build

build: bundle
	go build -ldflags "$(LDFLAGS)" -o bin/selfwatch

install: bundle
	go install -ldflags "$(LDFLAGS)"

test:
	go test -v ./...

release:
	gh release create "$$(date +%Y-%m-%d)" --generate-notes --title "Release $$(date +%Y-%m-%d)"
