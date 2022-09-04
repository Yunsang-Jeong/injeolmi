.PHONY: build clean

GIT_ROOT := $(shell git rev-parse --show-toplevel)

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o $(GIT_ROOT)/bin/main $(GIT_ROOT)/main.go

clean:
	rm -rf ./bin