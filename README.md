# Miracle74 API

Public API for accessing character data from **Miracle74.com**.  
Built in Go. OpenAPI-first.

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


