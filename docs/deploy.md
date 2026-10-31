# Deploy Sieve

## Local demo

```bash
go run ./cmd/sieve serve -config config.demo.yaml
# open http://localhost:8080/
curl -sf localhost:8080/healthz
curl -sf localhost:8080/v1/examples | head
```

## Railway

1. Create a Railway project/service from this repo (Dockerfile build via `railway.toml`).
2. The image entrypoint runs `sieve serve -config /config.demo.yaml` on `:8080`.
3. Health check path: `/healthz`.
4. Attach a public domain (Railway-generated or custom).

```bash
railway init   # or link an existing project
railway up
railway domain
```

### Ops note (2026-10-31)

Creating a new Railway project under the free plan may fail with
`Free plan resource provision limit exceeded`. Upgrade the Railway plan or free
capacity before `railway init` / `railway up`. Until then, run the demo locally
with `make run`.

## Custom domain `sieve.rejisterjack.com`

1. In Railway → service → Settings → Networking → Custom Domain, add `sieve.rejisterjack.com`.
2. At your DNS host, create the record Railway shows (usually CNAME to `*.up.railway.app`, or A/ALIAS as instructed).
3. Wait for TLS provisioning.
4. Verify:

   ```bash
   curl -sf https://sieve.rejisterjack.com/healthz
   curl -sf https://sieve.rejisterjack.com/ | head
   ```

Both should return HTTP 200.
