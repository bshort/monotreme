#!/bin/bash

# Rebuild and run Monotreme in Docker
# This script builds the frontend, builds the Docker image, and restarts the containers

set -e  # Exit on error

echo "=== Step 1: Building frontend ==="
cd frontend/web
npm run build
cd ../..

echo ""
echo "=== Step 2: Building Docker image ==="
docker build -f Dockerfile.local -t monotreme:local .

echo ""
echo "=== Step 3: Stopping existing containers ==="
docker compose down

echo ""
echo "=== Step 4: Starting containers ==="
docker compose up -d

echo ""
echo "=== Done! ==="
echo "Monotreme is now running at http://localhost:5231"
echo "Check logs with: docker compose logs -f"
