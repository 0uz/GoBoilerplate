# Go Boilerplate

A robust boilerplate with authentication built with Go, featuring JWT authentication, PostgreSQL integration, Redis caching, and RESTful API endpoints. This project follows clean architecture principles with a well-organized project structure.

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
- Postman collection for API testing

## Prerequisites

- Go 1.23.4 or higher
- PostgreSQL
- Valkey
- Docker and Docker Compose (optional)

## Setup and Installation
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
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_PASSWORD=  # if required
VALKEY_DB=0      # default database

# JWT Configuration
JWT_SECRET=s3cr3t
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=5h
```
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

Import the provided Postman collection (`boilerplate.postman_collection.json`) to test all available endpoints.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
