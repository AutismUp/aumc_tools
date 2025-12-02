#!/bin/bash

# Docker-based test runner for AUMC
# This script builds and runs the test environment using Docker

set -e

echo "========================================="
echo "AUMC Docker Test Environment"
echo "========================================="
echo ""

# Build the Go binary for Linux
echo "Building Linux binary..."
make build-linux
echo ""

# Build the Docker image
echo "Building Docker image..."
docker-compose build
echo ""

# Start the container
echo "Starting test container..."
docker-compose up -d
echo ""

# Wait for container to be ready
echo "Waiting for container to be ready..."
sleep 3

# Run the test script inside the container
echo "Running tests inside container..."
echo ""
docker-compose exec -T minecraft-test bash /workspace/test-aumc.sh

# Capture exit code
TEST_EXIT_CODE=$?

# Show how to access the container
echo ""
echo "========================================="
echo "Container is running!"
echo "========================================="
echo ""
echo "To access the container shell:"
echo "  docker-compose exec minecraft-test bash"
echo ""
echo "To stop the container:"
echo "  docker-compose down"
echo ""
echo "To view logs:"
echo "  docker-compose logs -f"
echo ""

exit $TEST_EXIT_CODE
