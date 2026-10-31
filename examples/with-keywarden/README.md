# Sieve + Keywarden perimeter

```text
untrusted content  →  Sieve (what agents may be told)
                   →  model
                   →  Keywarden (what agents may do / which secrets)
```

| Layer | Role |
|-------|------|
| **Sieve** | Scan user messages, scraped pages, tool outputs for prompt injection; block or flag before the model call. |
| **Keywarden** | Vault credentials, scoped task tokens, policy, and audited tool proxy — after the model decides to act. |

## Run

```bash
go run ./examples/with-keywarden
AGENT_INPUT='Summarize the harbor guide.' go run ./examples/with-keywarden
```

This example only depends on Sieve. Wire Keywarden using its task-token + MCP proxy docs when you need egress control.
