# Miracle74 API

High-performance REST API for accessing game data from [Miracle74.com](https://miracle74.com).

**Production API:** https://miracle74-api.fly.dev
**Documentation:** https://miracle74-api-docs.fly.dev

Built with Go, Redis caching, and OpenAPI-first design.

---

## Features

- **Character Data** - Get detailed character information including stats, guild, and death history
- **Power Gamers** - Top players ranked by level progression
- **Insomniacs** - Players ranked by online time
- **Guild Info** - Complete guild member listings with ranks and status
- **Redis Caching** - Fast response times with Valkey/Redis cache
- **Rate Limiting** - Built-in protection against abuse
- **OpenAPI Spec** - Full API documentation with Swagger UI

---

## API Endpoints

### Character
```bash
GET /characters/{name}
```
Get character information including level, vocation, guild, deaths, and more.

**Example:**
```bash
curl https://miracle74-api.fly.dev/characters/Oten
```

### Power Gamers
```bash
GET /powergamers?include_all=false
```
Get today's top power gamers (most levels gained).

**Example:**
```bash
curl https://miracle74-api.fly.dev/powergamers
```

### Insomniacs
```bash
GET /insomniacs?include_all=false
```
Get players with most online time.

**Example:**
```bash
curl https://miracle74-api.fly.dev/insomniacs
```

### Guilds
```bash
GET /guilds/{guildId}
```
Get complete guild information with all members.

**Example:**
```bash
curl https://miracle74-api.fly.dev/guilds/386
```

### Health Check
```bash
GET /health
```
Service health status.

---

## Documentation

Interactive API documentation with request/response examples:

**Swagger UI:** https://miracle74-api-docs.fly.dev

---

## Development

### Requirements

- Go 1.21+
- Docker (for Valkey/Redis)
- `mise` (optional, but recommended)

### Quick Start

```bash
git clone https://github.com/ethaan/miracle74-api.git
cd miracle74-api

# Start Redis cache
docker-compose up -d

# Run API server
go run cmd/api/main.go
```

Or with `mise`:

```bash
mise run docker-up
mise run dev
```

API runs at `http://localhost:8080`

### Run Documentation Server

```bash
go run cmd/docs/main.go
# or
mise run docs
```

Swagger UI at `http://localhost:8081`

### Available Tasks

```bash
mise run generate    # Generate API code from OpenAPI spec
mise run build       # Build binary
mise run test        # Run tests
mise run tidy        # Tidy Go modules
```

---

## Deployment

### API Server
```bash
fly deploy
```

### Documentation Server
```bash
fly deploy --config fly.docs.toml
```


