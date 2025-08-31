# Go Boilerplate

A robust boilerplate with authentication built with Go, featuring JWT authentication, PostgreSQL integration, Valkey caching, comprehensive monitoring, and RESTful API endpoints. This project follows clean architecture principles with a well-organized project structure.

## Features

- JWT-based authentication system
- Clean Architecture implementation with clear separation of concerns
- User management (registration, login, profile)
- Role-based access control (RBAC)
- Multi-device session management
- PostgreSQL database integration with GORM
- Valkey caching for improved performance
- RESTful API endpoints with middleware support
- Structured logging with logrus
- Docker and Docker Compose support
- Comprehensive monitoring with Prometheus & Grafana
- Real-time metrics and alerting
- Cross-platform build automation with Makefile
- Hoppscotch collections for API testing

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- Make (for build automation)

## Quick Start

```bash
# Start everything (app + monitoring)
make start-all

# Check status
make status
```

## Setup and Installation

### Using Docker (Recommended)
```bash
# Start the application
make run

# Start monitoring stack (optional)
make monitoring-start

# Or start everything at once
make start-all
```

### Manual Setup
1. Set up environment variables (create `.env` file):
```bash
# Server Configuration
PORT=8080
V1_PREFIX=/api/v1

# PostgreSQL Configuration
PG_DB_HOST=goauthboiler-postgres
PG_DB_PORT=5432
PG_DB_USER=app
PG_DB_PASSWORD=123456
PG_DB_NAME=appdb

# Valkey Configuration
VALKEY_HOST=goauthboiler-valkey
VALKEY_PORT=6379
VALKEY_PASSWORD=
VALKEY_DB=0

# JWT Configuration
JWT_SECRET=s3cr3t
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=5h
```

2. Start dependencies and run the application:
```bash
# Start dependencies (PostgreSQL & Valkey)
make dev

# Run the application locally
go run cmd/main.go
```

## Available Commands

Run `make help` to see all available commands:

```bash
# Application
make run              # Start the application
make stop             # Stop the application
make restart          # Restart the application
make rebuild          # Rebuild and restart

# Development
make dev              # Start dependencies only
make test             # Run tests
make logs             # Show application logs

# Monitoring
make monitoring-start # Start monitoring stack
make monitoring-stop  # Stop monitoring stack
make start-all        # Start app + monitoring
make stop-all         # Stop everything

# Utilities
make health           # Check service health
make clean            # Clean up containers
```

## Monitoring (Optional)

The project includes a comprehensive monitoring stack:

- **Grafana Dashboards**: http://localhost:3000 (admin/admin123)
- **Prometheus Metrics**: http://localhost:9090
- **Application Metrics**: http://localhost:8080/metrics

## API Documentation

### Authentication
All endpoints require `x-client-key` header

#### User Management
- `POST /users/registration` - Register a new user
- `POST /users/login` - Login and receive JWT tokens
- `GET /users/me` - Get current user profile

#### Authentication Management
- `POST /auth/token/refresh` - Refresh JWT token
- `POST /auth/logout` - Logout current session
- `POST /auth/logout/all` - Logout from all devices

#### Health Check
- `GET /live` - Liveness probe
- `GET /ready` - Readiness probe (checks database connection)

## Testing

### Unit Tests
```bash
make test             # Run all tests
make test-coverage    # Run tests with coverage
```

### API Testing
Import the provided test collections:
- **Hoppscotch Collection**: `doc/` 

## Development Workflow

```bash
# Development setup
make dev              # Start dependencies only
go run cmd/main.go    # Run app locally

# Testing
make test             # Run tests
make test-coverage    # Run with coverage

# Production
make start-all        # Start everything
make health           # Check health
make logs             # View logs

# Cleanup
make stop-all         # Stop everything
make clean            # Clean containers
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
