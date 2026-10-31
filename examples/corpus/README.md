# Injection corpus

Synthetic fixtures for regression tests and the Railway demo. Not harvested from live attacks.

| File | Purpose |
|------|---------|
| `safe.jsonl` | Inputs that must remain `safe` |
| `blocked.jsonl` | Inputs that must classify as `blocked` |
| `suspicious.jsonl` | Inputs that must classify as `suspicious` |
| `scraped_page.html` | Demo narrative: scraped page with hidden instructions |

Each JSONL line: `{"id","input","expect"}`.
