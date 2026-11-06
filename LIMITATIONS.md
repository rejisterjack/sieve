# Limitations

Sieve is a **defence-in-depth** layer, not a complete prompt-injection solution.

## What it catches

Known and near-known injection *patterns*: role overrides, ignore-previous-instruction phrasing, jailbreak keywords, delimiter/prompt smuggling markers, common obfuscation hints, and similar string-level signals. Built-in rules are regular expressions plus light normalization (zero-width stripping, whitespace collapse).

## What it misses

Contextual and persuasive attacks that do **not** match known patterns can get through. An attacker who writes novel instructions in natural language, or who hides intent across multiple benign-looking turns, may evade the filter.

Treat Sieve as one control among several (model-side policies, least-privilege tools, human review), not as a guarantee.

**Won't-build (closed):** ML / embedding-based / contextual multi-turn detection stays out of scope. Sieve remains pattern/regex defence-in-depth.

## False positives

Legitimate user messages can look like injections (“ignore the previous paragraph and summarize the next one”). Tune `block_on`, add domain-specific allow/deny rules in YAML, and review the structured flag log (`sieve logs`) before tightening thresholds in production.

## Pattern library drift

Detection quality tracks the pattern library. New techniques appear regularly; see [docs/patterns.md](docs/patterns.md) for how to extend rules and keep the corpus green.

## Companion, not replacement

[Keywarden](https://github.com/rejisterjack/keywarden) controls what agents can *do* (tools, credentials, policy). Sieve controls what agents can be *told*. Use both for a fuller perimeter; neither replaces the other.
