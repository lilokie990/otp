#!/bin/bash

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in your PATH. Please install Docker."
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    echo "Error: Docker Compose is not installed or not in your PATH. Please install Docker Compose."
    exit 1
fi

# Function to display usage information
function show_usage {
    echo "OTP Authentication System Runner"
    echo ""
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  start        Start the application (default if no command provided)"
    echo "  stop         Stop the application"
    echo "  logs         Show application logs"
    echo "  env          Show current environment variables"
    echo "  edit-env     Edit environment variables in docker-compose.yml"
    echo "  clean        Stop the application and remove all data"
    echo "  help         Show this help message"
    echo ""
}

# Process commands
case "$1" in
    start|"")
        echo "Starting OTP Authentication System..."
        docker compose up -d
        echo ""
        echo "Application is now running!"
        echo "- Web interface: http://localhost:8080"
        echo "- API Documentation: http://localhost:8080/swagger/index.html"
        echo "- Health Check: http://localhost:8080/health"
        ;;
    stop)
        echo "Stopping OTP Authentication System..."
        docker compose down
        echo "Application stopped."
        ;;
    logs)
        echo "Showing application logs (press Ctrl+C to exit)..."
        docker compose logs -f app
        ;;
    env)
        echo "Current environment variables in docker-compose.yml:"
        echo ""
        grep -A 20 "environment:" docker-compose.yml | grep -v "environment:" | grep -v "restart:" | grep -v -- "--" | sed 's/^[ \t]*//'
        ;;
    edit-env)
        if [[ "$EDITOR" ]]; then
            $EDITOR docker-compose.yml
        else
            echo "No default editor set. Please open docker-compose.yml manually."
            echo "Location: $(pwd)/docker-compose.yml"
        fi
        ;;
    clean)
        echo "Stopping application and removing all data..."
        docker compose down -v
        echo "Cleanup complete."
        ;;
    help)
        show_usage
        ;;
    *)
        echo "Error: Unknown command '$1'"
        echo ""
        show_usage
        exit 1
        ;;
esac

exit 0
