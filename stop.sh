#!/bin/bash

# Shutdown Monotreme Docker containers
# This script stops and removes all running Monotreme containers

set -e  # Exit on error

echo "=== Shutting down Monotreme ==="
docker compose down

echo ""
echo "=== Done! ==="
echo "All Monotreme containers have been stopped and removed."
echo "To start again, run: ./rebuild_and_run.sh"
