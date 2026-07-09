BIN=bin/server

run:
	go run ./cmd/server

build:
	go build -o $(BIN) ./cmd/server
