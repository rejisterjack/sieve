# Pattern library

Sieve ships built-in regex rules in `builtinRules()` (`sieve.go`). Rules evolve as new injection techniques appear.

## Adding a rule

1. Add a fixture line to `examples/corpus/{safe,blocked,suspicious}.jsonl` with the expected risk.
2. Copy the same JSONL (and HTML if needed) into `internal/httpapi/demo/corpus/` so the demo embed stays in sync.
3. Add or adjust a `Rule` in `builtinRules()` with a stable `Category` name.
4. Run `go test ./...` — corpus tests must pass with no safe→blocked regressions.
5. Document non-obvious patterns here if the name is not self-explanatory.
6. Keep [`LIMITATIONS.md`](../LIMITATIONS.md) honest: regex will miss contextual attacks.

## Normalization

Inputs are normalized before matching (`normalize.go`): lowercase, collapsed whitespace, zero-width characters stripped. Homoglyph attacks that survive Unicode lowercase may still need dedicated rules.

## Categories

Call `sieve.BuiltinCategoryCount()` or inspect `builtinRules()` for the current set (target ≥ 15).
