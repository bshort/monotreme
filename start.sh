#!/bin/bash

# Run Monotreme in Docker (LOCAL/DEVELOPMENT)
# This script builds and starts the containers with bridge networking
# For production deployment with MacVlan, use start-prod.sh

set -e  # Exit on error

echo "=== Step 1: Stopping existing containers ==="
docker compose down

echo ""
echo "=== Step 2: Building Docker image ==="
docker build -t monotreme:local -f Dockerfile .

echo ""
echo "=== Step 3: Starting containers ==="
docker compose up -d

echo ""
echo "=== Done! ==="
echo "Monotreme is now running at http://localhost:5231"
echo "Check logs with: docker compose logs -f"
