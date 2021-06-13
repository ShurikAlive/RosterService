all: build test check

modules:
	go mod tidy

build: modules
	go build cmd/*

test:
	go test ./...

check:
	golangci-lint run