# Sieve: Prompt Injection Filter for AI Agents

Private doc for me, Rupam. Not for public. Product name: Sieve. Home: sieve.rejisterjack.com. Go module under github.com/rejisterjack/sieve.

## The one-line idea

A Go middleware that screens every input going into an AI agent for prompt injection patterns, classifies the risk level, and blocks or flags malicious attempts before they can hijack the agent's behavior.

## The problem it solves

Useful agents read untrusted content constantly — user messages, web pages scraped during a task, tool outputs, uploaded files. Any of those can contain hidden instructions trying to override the agent's system prompt and make it do something it was never meant to do. This is prompt injection, and it is still one of the hardest unsolved security problems for AI agents in 2026.

Filtering suspicious strings is not enough once attacks become contextual and persuasive. Enterprise security tools exist but are complex, hosted, and priced for large teams. There is nothing lightweight, self-hosted, and open source. Sieve is that missing middleware layer. It pairs directly with Keywarden to form a complete security perimeter around an agent.

## What it does

- Middleware that intercepts every input going into an agent before it reaches the model
- Scans for known injection patterns: role override attempts, ignore-previous-instructions patterns, jailbreak templates, and contextual manipulation
- Classifies each input with a risk score: safe, suspicious, or blocked
- Suspicious inputs are flagged and logged; blocked inputs are rejected with a structured error before the API call is made
- Full log of every flagged input with the pattern that triggered it, for review and tuning
- Configurable rules: add custom patterns for your domain, adjust thresholds, set which risk levels to block vs. allow through
- Go library for embedding directly in an agent, or a standalone HTTP middleware for any language

## Who it is for

- Teams building agents that process untrusted content: user messages, scraped web pages, uploaded files, or tool outputs from third-party services
- Security-conscious developers who need an audit trail of suspicious inputs for compliance
- Teams already using Keywarden who want a complete security perimeter covering both input and output

## How it is used

1. Add Sieve as a Go dependency or run it as a sidecar HTTP service
2. Pipe every agent input through Sieve before it reaches the model
3. Sieve returns the input unchanged if it passes, or a structured error if it is blocked
4. Review flagged inputs in the log to tune thresholds for your use case
5. Add custom rules via a YAML config file for domain-specific patterns your agents are likely to see

## How I portray it

- One-line pitch: "Agents trust what they read. Sieve makes sure what they read is trustworthy."
- Lead with the demo: show a real prompt injection attempt — hidden instructions in a web page the agent scraped — being caught and blocked before it reaches the model
- Show the risk log because it makes the threat visible and concrete
- Position it as the natural companion to Keywarden: Keywarden controls what agents can do, Sieve controls what agents can be told

## How it runs

- Go library embedded directly in the agent's codebase, or run as a standalone HTTP sidecar service
- No hosted component; everything runs where the agent runs
- Config loaded from a YAML file for custom rules and thresholds
- Demo on Railway scans publicly available prompt injection examples to show it working

## Tech stack

- Go for the detection engine and the HTTP sidecar
- YAML config for rules and thresholds
- No database needed; flagged inputs are written to a structured log file
- Optional webhook for real-time alerts on blocked inputs

## What done looks like

- Go library that wraps any string input and returns a risk classification in one function call
- At least ten known injection pattern categories covered out of the box
- Custom rule support working via YAML config
- Structured log of every flagged input with the matched pattern, risk score, and timestamp
- HTTP sidecar mode working so any language can use Sieve over HTTP
- Benchmark showing how much latency the scan adds per call
- Demo GIF showing a real injection attempt being caught
- Repo at github.com/rejisterjack/sieve

## Honest notes

- No detection system catches everything; contextual and persuasive injections that do not match known patterns will get through — document this limitation honestly and frame Sieve as a defence-in-depth layer, not a complete solution
- False positives are a real risk; legitimate user messages may match injection patterns; tune thresholds carefully and make the custom rule system easy to use
- Detection quality is only as good as the pattern library; plan to update it regularly as new injection techniques emerge
- Pairs with Keywarden but does not replace it; they protect different surfaces

## Deployment

- No hosting needed for the library mode
- Railway for the HTTP sidecar demo
- GitHub releases for prebuilt sidecar binaries
- go install for developers who have Go installed
- No Vercel, no Neon
