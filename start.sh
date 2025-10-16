#!/bin/bash

# Run Monotreme in Docker
# This script starts the containers without rebuilding

set -e  # Exit on error

echo "=== Step 1: Stopping existing containers ==="
docker compose down

echo ""
echo "=== Step 2: Starting containers ==="
docker compose up -d

echo ""
echo "=== Done! ==="
echo "Monotreme is now running at http://localhost:5231"
echo "Check logs with: docker compose logs -f"
