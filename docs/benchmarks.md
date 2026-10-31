# Benchmarks

Measure Scan latency with the Go toolchain:

```bash
make bench
# or
go test -bench=BenchmarkScan -benchmem ./...
```

## Sample results

Machine: Apple M5 (darwin/arm64), Go 1.27.3, 2026-10-31.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| BenchmarkScanSafe | ~27,900 | 68 | 1 |
| BenchmarkScanBlocked | ~37,300 | 235 | 3 |
| BenchmarkScanScrapedHTML | ~236,000 | 1,360 | 5 |

Numbers vary by CPU and Go version; re-run locally before citing in production capacity planning.
