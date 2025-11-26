.PHONY: run build test lint fmt

run:
	go run ./cmd/api

build:
	go build -o bin/wound_iq_api_new ./cmd/api

test:
	go test ./... -v

lint:
	@golangci-lint run || true

fmt:
	gofmt -s -w .
