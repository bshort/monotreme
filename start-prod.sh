#!/bin/bash

# Run Monotreme in Docker (PRODUCTION)
# This script builds and starts the containers with MacVlan networking
# For local development with bridge networking, use start.sh

set -e  # Exit on error

echo "=== Step 1: Stopping existing containers ==="
docker compose -f docker-compose.yml -f docker-compose.prod.yml down

echo ""
echo "=== Step 2: Building Docker image ==="
docker build -t monotreme:local -f Dockerfile .

echo ""
echo "=== Step 3: Starting containers ==="
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

echo ""
echo "=== Done! ==="
echo "Monotreme is now running with MacVlan networking"
echo "Check logs with: docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f"
