#!/bin/bash
# Exit immediately if a command exits with a non-zero status.
set -e

# Build the Go backend
echo "Building Go backend..."
go build -o bin/api cmd/api/main.go

# Build the frontend
echo "Building frontend..."
(cd web && npm install && npm run build)

echo "Build complete."
