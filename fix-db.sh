#!/bin/bash

echo "=========================================="
echo "Fixing Database Container"
echo "=========================================="

# Stop containers
echo "1. Stopping containers..."
docker compose -f docker-compose.prod.yml down

# Remove old containers
echo "2. Removing old containers..."
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f 2>/dev/null || true

# Kill port 5432
echo "3. Killing processes on port 5432..."
sudo fuser -k 5432/tcp 2>/dev/null || true

# Remove volumes (optional - uncomment if needed)
# echo "4. Removing volumes..."
# docker volume ls | grep tukem | awk '{print $2}' | xargs docker volume rm 2>/dev/null || true

# Start fresh
echo "5. Starting containers..."
docker compose -f docker-compose.prod.yml up -d

# Wait for database
echo "6. Waiting for database to be ready..."
sleep 15

# Check status
echo ""
echo "=========================================="
echo "Container Status:"
echo "=========================================="
docker compose -f docker-compose.prod.yml ps

echo ""
echo "=========================================="
echo "Database Logs (last 20 lines):"
echo "=========================================="
docker compose -f docker-compose.prod.yml logs db | tail -20

echo ""
echo "=========================================="
echo "Done!"
echo "=========================================="
echo ""
echo "If database is running, now run:"
echo "  ./setup-database.sh"
echo ""

