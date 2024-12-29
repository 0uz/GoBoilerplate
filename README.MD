# Go Auth Boilerplate

A robust authentication boilerplate built with Go, featuring JWT authentication, PostgreSQL integration, and RESTful API endpoints.

## Features

- JWT-based authentication
- User registration and login
- PostgreSQL database integration
- RESTful API endpoints
- Secure password hashing
- Environment variable configuration
- Middleware for route protection

## Prerequisites

- Go 1.16+
- PostgreSQL
- Docker (optional)

## Setup

1. Set up environment variables (create `.env` file):
```bash
PORT=8080
V1_PREFIX=/api/v1
PG_DB_HOST=localhost
PG_DB_PORT=5432
PG_DB_USER=goauthboiler
PG_DB_PASSWORD=123456
PG_DB_NAME=postgres

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=s3cr3t
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=5h

```

3. Install dependencies:
```bash
go mod download
```

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Health Check Endpoints
- `GET /live` - Liveness probe
- `GET /ready` - Readiness probe

### Authentication Endpoints
All endpoints require `x-client-key` header

#### User Registration
```bash
POST /api/v1/user/register
Content-Type: application/json

{
    "email": "example@test.com",
    "username": "example@test.com",
    "password": "yourpassword"
}
```

#### User Login
```bash
GET /api/v1/user/login
Content-Type: application/json

{
    "email": "example@test.com",
    "password": "yourpassword"
}
```

#### Anonymous User Creation
```bash
POST /api/v1/user/anonymous
```

#### Anonymous User Login
```bash
GET /api/v1/user/anonymous/login
Content-Type: application/json

{
    "email": "anonymous-uuid"
}
```

#### User Logout
```bash
POST /api/v1/user/logout
Authorization: Bearer <access_token>
```

#### Refresh Token
```bash
POST /api/v1/user/refresh
Content-Type: application/json

{
    "refresh_token": "your-refresh-token"
}
```

#### Confirm User
```bash
GET /api/v1/user/{user-id}/confirm
```

### Required Headers
- `x-client-key`: Client application identifier
- `Authorization`: Bearer token (for authenticated endpoints)
- `Content-Type`: application/json (for requests with body)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
