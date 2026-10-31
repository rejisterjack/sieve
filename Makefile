.PHONY: test vet build run bench logs

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -o bin/sieve ./cmd/sieve

run:
	go run ./cmd/sieve serve -config config.demo.yaml

bench:
	go test -bench=BenchmarkScan -benchmem ./...

logs:
	go run ./cmd/sieve logs -path sieve.log -summary
