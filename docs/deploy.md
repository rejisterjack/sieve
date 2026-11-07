# Deploy Sieve

## Local demo

```bash
go run ./cmd/sieve serve -config config.demo.yaml
# open http://localhost:8080/
curl -sf localhost:8080/healthz
curl -sf localhost:8080/v1/examples | head
```

## Render (public demo)

The live host [sieve.rejisterjack.com](https://sieve.rejisterjack.com) runs on Render
(`Dockerfile` → free Web Service). The process binds `0.0.0.0:$PORT` when Render
sets `PORT` (local default remains `:8080`). Health check path: `/healthz`.

1. Create a Render Web Service from this repo (Docker runtime, `./Dockerfile`).
2. Note the `*.onrender.com` URL (e.g. `https://sieve-ih2z.onrender.com`).
3. Either attach `sieve.rejisterjack.com` as a Render custom domain, **or** (when
   DNS for `rejisterjack.com` already lives on Vercel) deploy `edge/` as a
   Vercel rewrite to that Render URL and alias the custom domain there — same
   pattern as Tock.

```bash
curl -sf https://sieve.rejisterjack.com/healthz
curl -sf https://sieve.rejisterjack.com/ | head
```

Free Render instances sleep after ~15 minutes idle; the first request after
sleep can be slow.

## Railway (optional)

Railway also works for the same Dockerfile (`railway.toml`). Free-plan limits may
block new projects or custom domains; prefer Render for the public demo unless
you already have Railway capacity.
