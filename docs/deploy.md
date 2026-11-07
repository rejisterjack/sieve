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
2. The image entrypoint runs `sieve serve -config /config.demo.yaml`. The process
   binds `0.0.0.0:$PORT` when Railway sets `PORT` (local default remains `:8080`).
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

Railway’s free plan may refuse custom domains. When DNS for `rejisterjack.com`
already lives on Vercel (wildcard `*`), put a rewrite in front of the Railway
URL the same way Tock does:

1. Deploy the Go sidecar on Railway (Dockerfile / `railway.toml`) and note the
   `*.up.railway.app` service domain.
2. Point `edge/vercel.json` at that Railway URL, then deploy the `edge/`
   directory as a Vercel project (`rootDirectory=edge`).
3. Attach `sieve.rejisterjack.com` as an alias on that edge project (removes the
   dead `DEPLOYMENT_NOT_FOUND` wildcard hit).
4. Verify:

   ```bash
   curl -sf https://sieve.rejisterjack.com/healthz
   curl -sf https://sieve.rejisterjack.com/ | head
   ```

Both should return HTTP 200. The edge project only rewrites — it does not host
the Go binary.

If your Railway plan allows custom domains, you can skip the edge proxy and
attach `sieve.rejisterjack.com` directly in Railway networking instead.
