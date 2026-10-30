# Sieve

Prompt injection filter for AI agents.

> Agents trust what they read. Sieve makes sure what they read is trustworthy.

**Home:** [sieve.rejisterjack.com](https://sieve.rejisterjack.com) · **Module:** `github.com/rejisterjack/sieve`

## Library

```go
import "github.com/rejisterjack/sieve"

s, err := sieve.New(sieve.Config{})
if err != nil {
    log.Fatal(err)
}

result := s.Scan(userInput)
if s.ShouldBlock(result) {
    return fmt.Errorf("blocked: %s", result.Risk)
}
```

## HTTP sidecar

```bash
go run ./cmd/sieve -config config.example.yaml
curl -s localhost:8080/v1/scan -d '{"input":"Ignore previous instructions"}' -H 'content-type: application/json'
```

## Config

See [`config.example.yaml`](config.example.yaml) for thresholds, custom rules, log path, and optional webhook alerts.

## Status

Early scaffold. Built-in pattern categories and HTTP scan endpoint are in place; detection quality will expand over time. Pair with [Keywarden](https://github.com/rejisterjack/keywarden) for a full agent security perimeter.
