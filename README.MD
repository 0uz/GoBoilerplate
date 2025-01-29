# Go Auth Boilerplate

A robust authentication boilerplate built with Go, featuring JWT authentication, PostgreSQL integration, Redis caching, and RESTful API endpoints. This project follows clean architecture principles with a well-organized project structure.

## Features

- JWT-based authentication
- User registration and login
- Anonymous user support
- PostgreSQL database integration with GORM
- Redis cache integration for improved performance
- RESTful API endpoints
- Environment variable configuration
- Clean Architecture implementation
- Structured logging with logrus

## Prerequisites
- Go 1.23+
- PostgreSQL
- Redis
- Docker (optional)

## Setup
1. Set up environment variables (create `.env` file):
```bash
# Server Configuration
PORT=8080
V1_PREFIX=/api/v1

# PostgreSQL Configuration
PG_DB_HOST=localhost
PG_DB_PORT=5432
PG_DB_USER=goauthboiler
PG_DB_PASSWORD=123456
PG_DB_NAME=postgres

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=  # if required
REDIS_DB=0      # default database

# JWT Configuration
JWT_SECRET=s3cr3t
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=5h
```

3. Install dependencies:
```bash
go mod download
```

4. Run with Docker (optional):
```bash
docker-compose up -d
```

   Or run locally:
```bash
go run cmd/main.go
```

## API Endpoints

### Health Check Endpoints
- `GET /live` - Liveness probe
- `GET /ready` - Readiness probe (checks database connection)

### Authentication Endpoints
All endpoints require `x-client-key` header
You can import postman collection to test all endpoints
- `POST /users/registration` - Register a new user
- `POST /users/login` - Login and receive JWT tokens
- `POST /auth/token/refresh` - Refresh JWT token
- `POST /auth/logout` - Logout
- `POST /auth/logout/all` - Logout from all devices
- `GET /users/me` - Get current user


## License

This project is licensed under the MIT License - see the LICENSE file for details
