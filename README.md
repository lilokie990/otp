# OTP Authentication System - Documentation

A comprehensive OTP (One-Time Password) based authentication system implemented in Go using the Gin framework.

## Table of Contents

- [OTP Authentication System - Documentation](#otp-authentication-system---documentation)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [System Architecture](#system-architecture)
  - [Project Structure](#project-structure)
  - [Setup and Installation](#setup-and-installation)
    - [Prerequisites](#prerequisites)
    - [Running with Docker (Recommended)](#running-with-docker-recommended)
    - [Running Locally](#running-locally)
  - [Configuration](#configuration)
  - [Swagger Documentation](#swagger-documentation)
    - [Using Existing Swagger Documentation](#using-existing-swagger-documentation)
    - [Updating Swagger Documentation](#updating-swagger-documentation)
    - [Accessing Swagger UI](#accessing-swagger-ui)
    - [API Annotation Examples](#api-annotation-examples)
  - [API Reference](#api-reference)
    - [Authentication Endpoints](#authentication-endpoints)
    - [User Endpoints](#user-endpoints)
  - [Security Considerations](#security-considerations)
  - [Troubleshooting](#troubleshooting)
    - [Common Issues](#common-issues)
    - [Logs](#logs)

## Overview

This system provides a secure and reliable way to authenticate users using one-time passwords sent to their phone numbers. It utilizes modern web technologies and follows best practices for security and performance.

## Features

- **OTP-based Authentication**: Generate and verify one-time passwords for secure user authentication
- **JWT Token Management**: Secure API access using JSON Web Tokens
- **Rate Limiting**: Prevent abuse with configurable rate limiting
- **User Management**: Create and manage user accounts
- **API Documentation**: Interactive Swagger documentation
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Database Migrations**: Automatic database setup and migrations
- **Environment Configuration**: Flexible environment-based configuration

## System Architecture

The system is built using the following components:

- **Go with Gin Framework**: For the core HTTP server and routing
- **PostgreSQL**: For persistent data storage (user records)
- **Redis**: For OTP storage and rate limiting
- **Docker & Docker Compose**: For containerization and easy deployment
- **Swagger**: For API documentation

## Project Structure

```plaintext
├── cmd/                    # Application entry points
│   └── main.go             # Main application file
├── config/                 # Configuration handling
│   └── config.go           # Configuration structures and loaders
├── docs/                   # Documentation
│   └── swagger/            # Swagger API documentation
├── internal/               # Private application code
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models and DTOs
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic layer
│   └── utils/              # Utility functions
├── migrations/             # Database migrations
│   └── 001_create_users_table.sql
├── Dockerfile              # Docker build instructions
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

## Setup and Installation

### Prerequisites

- Docker and Docker Compose (for containerized setup)
- Go 1.24 or later (for local development)
- PostgreSQL (for local development)
- Redis (for local development)

### Running with Docker (Recommended)

1. Extract the ZIP file to a convenient location

2. Start the application using Docker Compose:

   ```bash
   # Start the application
   docker compose up -d
   ```

3. Access the application:
   - Web interface: `http://localhost:8080`
   - API Documentation: `http://localhost:8080/swagger/index.html`
   - Health Check: `http://localhost:8080/health`

4. To stop the application:

   ```bash
   # Stop the application
   docker compose down

   # To clean up data (optional)
   docker compose down -v
   ```

### Running Locally

1. Install Go 1.24 or later:

   ```bash
   # For macOS with Homebrew
   brew install go

   # For Linux
   # Download from https://golang.org/dl/ and follow installation instructions

   # Verify installation
   go version
   ```

2. Set up PostgreSQL and Redis:

   ```bash
   # For macOS with Homebrew
   brew install postgresql redis
   brew services start postgresql
   brew services start redis

   # For Linux
   sudo apt-get install postgresql redis-server
   sudo systemctl start postgresql
   sudo systemctl start redis
   ```

3. Clone the repository:

   ```bash
   git clone https://github.com/lilokie990/otp.git
   cd otp
   ```

4. Install dependencies:

   ```bash
   go mod download

   # Install Swagger tools for API documentation
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

5. Generate Swagger documentation:

   ```bash
   swag init -g cmd/main.go
   ```

6. Configure the application (see [Configuration](#configuration)):

   ```bash
   # Create local config file
   cp config.example.yaml config.local.yaml

   # Edit the config file with your settings
   vim config.local.yaml
   ```

7. Create and initialize database:

   ```bash
   # Connect to PostgreSQL
   psql -U postgres

   # In PostgreSQL shell
   CREATE DATABASE otpauth;
   \q

   # Apply migrations
   # If using a migration tool like golang-migrate
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/otpauth?sslmode=disable" up
   ```

8. Run the application:

   ```bash
   go run cmd/main.go
   ```

9. Test the application:

   ```bash
   # Request OTP
   curl -X POST http://localhost:8080/v1/auth/request-otp \
     -H "Content-Type: application/json" \
     -d '{"phone_number": "989123456789"}'

   # Check logs to see the OTP code
   # Verify OTP (replace 123456 with the actual OTP from logs)
   curl -X POST http://localhost:8080/v1/auth/verify-otp \
     -H "Content-Type: application/json" \
     -d '{"phone_number": "989123456789", "otp": "123456"}'
   ```

## Configuration

The application can be configured using a YAML file (`config.yaml`). Here's an example configuration:

```yaml
service:
  name: "otp-auth-service"
  env: "development"
  gracefulShutdownSecond: 5
  http:
    port: "8080"

postgres:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  databaseName: "otpauth"
  sslMode: "disable"
  timeZone: "UTC"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expirationHours: 24

otp:
  expiration: 120  # seconds
  length: 6
  rateLimit:
    count: 3
    time: 10  # minutes
```

You can override the configuration file path by setting the `CONFIG_PATH` environment variable. If the configuration file is not found, the application will fall back to environment variables with the same names as the configuration values.

## Swagger Documentation

This project uses Swagger/OpenAPI for API documentation. Swagger provides interactive documentation that allows you to explore and test API endpoints directly from a web interface.

### Using Existing Swagger Documentation

The project already includes pre-generated Swagger documentation in the `/docs` folder:

- `/docs/swagger/docs.go` - Generated API definitions
- `/docs/swagger/swagger.json` - OpenAPI spec in JSON format
- `/docs/swagger/swagger.yaml` - OpenAPI spec in YAML format

### Accessing Swagger UI

Once the application is running, you can access the Swagger UI at:

```plaintext
http://localhost:8080/swagger/index.html
```

### API Annotation Examples

Swagger uses annotations in your handler functions to generate documentation. Here's an example from the project:

```go
// RequestOTP handles OTP request
// @Summary Request OTP for a phone number
// @Description Generate and send a one-time password to the provided phone number (OTP is printed to server logs)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RequestOTPRequest true "Phone number to send OTP to"
// @Success 200 {object} models.RequestOTPResponse "OTP sent successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 429 {object} models.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
    // Handler implementation
}
```

## API Reference

### Authentication Endpoints

- **Request OTP**: `POST /v1/auth/request-otp`

  ```json
  {
    "phone_number": "+989123456789"
  }
  ```

  Response:

  ```json
  {
    "message": "OTP sent successfully. Check server logs for the code."
  }
  ```

  Accepted Iranian phone number formats:
  - International: `+989123456789`
  - National: `09123456789`
  - Short: `9123456789`

  The system validates Iranian mobile network prefixes including MCI (910-919, 990-996), Irancell (930-939, 901-905), and RighTel (920-922).

  **Note:** For security reasons, OTP codes are not included in the API response. Instead, they are printed to the server logs in the format: `[OTP] Phone: +989123456789, Code: 123456`- **Verify OTP**: `POST /v1/auth/verify-otp`

  ```json
  {
    "phone_number": "+989123456789",
    "otp": "123456"
  }
  ```

  OTP validation requirements:
  - Must be exactly 6 digits
  - Must be numeric only

### User Endpoints

All user endpoints require JWT authentication via the Authorization header.

- **Get User**: `GET /v1/users/:id`
  - Requires: Authorization header with Bearer token

- **List Users**: `GET /v1/users`
  - Requires: Authorization header with Bearer token
  - Query Parameters:
    - `page`: Page number (default: 1)
    - `pageSize`: Items per page (default: 10)
    - `search`: Search term for phone number

## Security Considerations

- OTPs expire after a configurable period (default: 120 seconds)
- Rate limiting prevents brute force attacks (default: 3 attempts per 10 minutes)
- JWT tokens expire after a configurable period (default: 24 hours)
- Database credentials should be securely managed in production
- Use HTTPS in production environments

## Troubleshooting

### Common Issues

1. **Connection to database failed**:
   - Verify that PostgreSQL is running
   - Check database credentials
   - Ensure network connectivity between app and database

2. **Connection to Redis failed**:
   - Verify that Redis is running
   - Check Redis credentials
   - Ensure network connectivity between app and Redis

3. **OTP not received**:
   - In development, OTP is printed to server logs (not included in API response)
   - Check server logs to see generated OTPs
   - Check rate limiting configuration
   - Verify phone number format

4. **JWT verification failed**:
   - Check JWT_SECRET environment variable
   - Verify token expiration
   - Ensure correct token format (Bearer token)

### Logs

To view application logs:

```bash
docker compose logs -f app
```

For specific logs:

```bash
docker compose logs -f postgres  # Database logs
docker compose logs -f redis     # Redis logs
```

