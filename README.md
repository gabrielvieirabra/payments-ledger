# payments-ledger

Payments ledger service built with Go.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- [golangci-lint](https://golangci-lint.run/) (optional, for linting)

## Getting Started

```bash
# Copy environment variables
cp .env.example .env

# Run locally
make run

# Run with Docker
make docker-up

# Run tests
make test
```

## Project Structure

```
.
├── cmd/api/          # Application entrypoint
├── internal/
│   ├── config/       # Configuration loading
│   ├── domain/       # Domain entities and interfaces
│   ├── handler/      # HTTP handlers
│   ├── repository/   # Data access layer
│   └── service/      # Business logic
├── pkg/              # Shared libraries (exported)
├── api/              # API specs (OpenAPI, protobuf)
├── migrations/       # Database migrations
├── scripts/          # Build and automation scripts
├── docs/             # Documentation
├── Dockerfile        # Multi-stage production build
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Available Make Targets

```bash
make help
```

## Environment Variables

See [.env.example](.env.example) for all available configuration options.