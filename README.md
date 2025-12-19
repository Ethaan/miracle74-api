# Miracle74 API

REST API for [Miracle74.com](https://miracle74.com) game data.

**API:** https://miracle74-api.fly.dev
**Docs:** https://miracle74-api-docs.fly.dev

---

## Development

```bash
git clone https://github.com/ethaan/miracle74-api.git
cd miracle74-api

# Start cache and run API
docker-compose up -d
go run cmd/api/main.go
```

Or with `mise`:
```bash
mise run docker-up
mise run dev
```

API: `http://localhost:8080`
Docs: `mise run docs` → `http://localhost:8081`

---

## Deployment

```bash
fly deploy                        # API
fly deploy --config fly.docs.toml # Docs
```
