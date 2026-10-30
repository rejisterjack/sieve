.PHONY: test vet build run

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -o bin/sieve ./cmd/sieve

run:
	go run ./cmd/sieve -config config.example.yaml
