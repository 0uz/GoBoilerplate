# Go Boilerplate

Production-ready Go boilerplate with JWT authentication, clean architecture, and full observability stack.

## Features

- **Authentication**: JWT-based auth with access/refresh tokens, multi-device session management
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Database**: PostgreSQL with GORM, connection pooling, migrations
- **Caching**: Valkey (Redis-compatible) for sessions and caching
- **Observability**: OpenTelemetry integration (traces, metrics, logs), Prometheus, Grafana, Loki, Tempo
- **Configuration**: YAML-based config with environment overrides (Viper)
- **Security**: Role-based access control, rate limiting, client authentication

## Quick Start

```bash
# Clone and setup
git clone https://github.com/ouz/goauthboilerplate.git
cd goauthboilerplate

# Start with Docker
make run

# Or start dependencies only for local development
make dev
go run cmd/main.go
```

## Configuration

Configuration uses YAML files with environment-based overrides:

```
config/
├── base.yaml              # Base configuration
├── development.yaml       # Local development overrides
├── development-docker.yaml # Docker development overrides
└── production.yaml        # Production overrides
```

Set `APP_ENV` to select environment (e.g., `APP_ENV=production`).

Sensitive values can be overridden via environment variables:
```bash
export POSTGRES_PASSWORD=your-secure-password
export JWT_SECRET=your-32-char-minimum-secret
```

## Project Structure

```
├── cmd/                    # Application entrypoint
├── config/                 # YAML configuration files
├── internal/
│   ├── adapters/          # HTTP handlers, repositories, middleware
│   ├── application/       # Business logic services
│   ├── config/            # App configuration loader
│   ├── domain/            # Domain models and interfaces
│   └── observability/     # Telemetry setup
├── pkg/                   # Reusable packages
│   ├── auth/              # Auth utilities (token types, client types)
│   ├── cache/             # Cache interfaces and implementations
│   ├── config/            # Generic config loader
│   ├── errors/            # Error handling
│   ├── log/               # Structured logging
│   ├── otel/              # OpenTelemetry helpers
│   ├── response/          # HTTP response utilities
│   └── stream/            # Event streaming (Redis Streams)
├── migrations/            # Database migrations
└── monitoring/            # Grafana, Prometheus, Loki, Tempo configs
```

## API Endpoints

All business endpoints require `x-client-key` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users/registration` | Register user |
| POST | `/api/v1/users/login` | Login |
| GET | `/api/v1/users/me` | Get current user (auth required) |
| POST | `/api/v1/auth/token/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout current session |
| POST | `/api/v1/auth/logout/all` | Logout all sessions |
| GET | `/live` | Liveness probe |
| GET | `/ready` | Readiness probe |
| GET | `/metrics` | Prometheus metrics |

## Commands

```bash
make run              # Start application
make dev              # Start dependencies only
make stop             # Stop application
make test             # Run tests
make monitoring-start # Start observability stack
make start-all        # Start everything
make status           # Show service status
make logs             # View logs
```

## Monitoring

Start the monitoring stack with `make monitoring-start`:

- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **App Metrics**: http://localhost:8080/metrics

## Requirements

- Go 1.25+
- Docker & Docker Compose
- Make

## License

MIT License - see [LICENSE](LICENSE) for details.
