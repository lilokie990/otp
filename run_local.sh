#!/bin/bash

echo "Starting OTP Authentication System in local development mode..."

# Set environment variables for local development
export DB_HOST=localhost
export REDIS_HOST=localhost
# Use a different port to avoid conflicts
export PORT=8082

# Set config path to use local configuration
export CONFIG_PATH="$(pwd)/config.local.yaml"

# Make sure PostgreSQL and Redis are running
echo "Checking if PostgreSQL is running..."
pg_isready -h localhost -p 5432 || echo "WARNING: PostgreSQL might not be running"

echo "Checking if Redis is running..."
redis-cli ping > /dev/null 2>&1 || echo "WARNING: Redis might not be running"

# Run the application
echo "Starting application with local config..."
go run cmd/main.go
