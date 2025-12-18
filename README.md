# Miracle74 API

Public API for accessing character data from **Miracle74.com**.
Built in Go. OpenAPI-first.

---

## Deploy to Fly.io

[![Deploy on Fly.io](https://fly.io/button.svg)](https://fly.io/launch?template=https://github.com/ethaan/miracle74-api)

### Setup Instructions

1. Click the deploy button above (requires a Fly.io account)
2. Set up Upstash Redis:
   - Create a free account at [Upstash](https://upstash.com/)
   - Create a new Redis database (select the Fly.io region closest to your app)
   - Copy the **Redis Connect URL** (format: `redis://default:password@host.upstash.io:6379`)
3. Set the `CACHE_URL` secret in Fly.io:
   ```bash
   fly secrets set CACHE_URL="redis://default:YOUR_PASSWORD@your-host.upstash.io:6379"
   ```

The cache client automatically detects Upstash URLs and enables TLS.

---

## Requirements

- Go 1.21+
- Docker (for Valkey/Redis)
- `mise` (optional)

---

## Run (Dev)

```bash
git clone https://github.com/ethaan/miracle74-api.git
cd miracle74-api

docker-compose up -d
go run cmd/api/main.go
```

Or with `mise`:

```bash
mise run docker-up
mise run dev
```

API runs at:

```
http://localhost:8080
```

---

## Swagger / OpenAPI

```bash
go run cmd/docs/main.go
```

Or:

```bash
mise run docs
```

Swagger UI:

```
http://localhost:8081
```

OpenAPI spec:

```
openapi.yaml
```


