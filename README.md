# Miracle74 API

A high-performance Go API for accessing character data from Miracle74.com. Features OpenAPI-first design, Valkey caching, and comprehensive character information extraction.

## Features

- RESTful API with OpenAPI 3.0 specification
- Character data extraction including stats, guild info, and death history
- Valkey/Redis caching with configurable TTL
- Rate limiting support
- Interactive API documentation (Swagger UI)
- Hot reload development environment

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for Valkey)
- mise (optional, for task running)

### Installation

```bash
git clone https://github.com/ethaan/miracle74-api.git
cd miracle74-api
go mod download
```

### Running the API

Start Valkey cache:
```bash
docker-compose up -d
# or with mise:
mise run docker-up
```

Run the API server:
```bash
go run cmd/api/main.go
# or with mise and hot reload:
mise run dev
```

The API will be available at `http://localhost:8080`

### API Documentation

View interactive documentation:
```bash
go run cmd/docs/main.go
# or with mise:
mise run docs
```

Open `http://localhost:8081` in your browser to access Swagger UI.

## API Endpoints

### Get Character Information

```
GET /characters/{name}
```

Returns comprehensive character data including:
- Basic info (name, sex, vocation, level)
- Guild membership and rank
- Last login timestamp
- Premium account status
- Recent death history

Example:
```bash
curl http://localhost:8080/characters/Oten
```

### Health Check

```
GET /health
```

Returns API health status.

## Development

### Project Structure

```
miracle74-api/
├── cmd/
│   ├── api/          # Main API server
│   └── docs/         # Documentation server
├── internal/
│   ├── api/          # Generated OpenAPI code (do not edit)
│   ├── handlers/     # HTTP handlers
│   ├── services/     # Business logic
│   ├── repo/         # Data repositories
│   └── types/        # Shared types
├── pkg/
│   ├── cache/        # Valkey cache client
│   └── miracle74/    # Web scraper client
├── openapi.yaml      # API specification
└── docker-compose.yml
```

### Making Changes

1. Update `openapi.yaml` to modify the API contract
2. Regenerate code: `mise run generate`
3. Implement handlers in `internal/handlers/`
4. Test with hot reload: `mise run dev`

### Available Tasks

```bash
mise run generate   # Regenerate OpenAPI code
mise run dev        # Run API with hot reload
mise run docs       # Run documentation server
mise run build      # Build production binary
mise run test       # Run tests
mise run docker-up  # Start Valkey
mise run docker-down # Stop Valkey
```

### Environment Variables

- `PORT` - API server port (default: 8080)
- `CACHE_URL` - Valkey connection string (default: localhost:6379)
- `DOCS_PORT` - Documentation server port (default: 8081)

## Architecture

### OpenAPI-First Design

This project follows an OpenAPI-first approach:

1. Define API contract in `openapi.yaml`
2. Generate server code with ogen
3. Implement business logic in handlers and services
4. Changes to the API require updating the spec first

### Caching Strategy

- Character data is cached by name with a 10-minute TTL
- Cache keys: `character:{name}`
- Cache misses trigger web scraping from miracle74.com
- Scraped HTML is saved to `public/` directory during development

### Data Flow

```
HTTP Request → Handler → Service → Repository → Cache
                           ↓
                     Miracle74 Client
```

## Contributing

Contributions are welcome. Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes following existing code style
4. Update tests if applicable
5. Update OpenAPI spec if modifying API contracts
6. Submit a pull request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Add comments for exported functions

## Testing

```bash
go test -v ./...
# or with mise:
mise run test
```

## Building for Production

```bash
go build -o bin/api cmd/api/main.go
# or with mise:
mise run build
```

## License

This project is provided as-is for educational and research purposes.

## Acknowledgments

Built for the Miracle74 community. Visit [miracle74.com](https://miracle74.com) for more information about the game.
