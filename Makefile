.PHONY: generate test build build-cli build-server install

generate:
	buf generate

test: generate
	go test ./...

build-cli: generate
	go build -o dj ./cli

build-server: generate
	go build -o dj-server ./backend/cmd/server

build: build-cli build-server

install: build-cli
	go install ./cli
