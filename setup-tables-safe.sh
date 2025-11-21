#!/bin/bash

# Script untuk setup tables dengan handling password yang mengandung karakter khusus
# Usage: ./setup-tables-safe.sh

set -e

echo "=========================================="
echo "Setting up Database Tables (Safe Mode)"
echo "=========================================="

# Check if container is running
if ! docker ps | grep -q "tukem-db-prod"; then
    echo "Error: Database container is not running!"
    echo "Please start it first: docker compose -f docker-compose.prod.yml up -d db"
    exit 1
fi

# Load .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Use PGPASSWORD to avoid password issues
export PGPASSWORD="${DB_PASSWORD}"

echo "Waiting for database to be ready..."
for i in {1..60}; do
    # Try to connect using PGPASSWORD
    if docker compose -f docker-compose.prod.yml exec -T -e PGPASSWORD="${DB_PASSWORD}" tukem-db-prod psql -U "${DB_USER}" -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    if [ $((i % 5)) -eq 0 ]; then
        echo "Waiting... ($i/60)"
    fi
    sleep 1
done

# Check if database is ready
if ! docker compose -f docker-compose.prod.yml exec -T -e PGPASSWORD="${DB_PASSWORD}" tukem-db-prod psql -U "${DB_USER}" -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
    echo ""
    echo "Error: Cannot connect to database!"
    echo ""
    echo "Troubleshooting:"
    echo "1. Check container status:"
    docker ps -a | grep tukem-db-prod
    echo ""
    echo "2. Check database logs:"
    docker logs tukem-db-prod 2>&1 | tail -20
    echo ""
    exit 1
fi

# Run init script
echo ""
echo "Running database initialization..."
docker compose -f docker-compose.prod.yml exec -T -e PGPASSWORD="${DB_PASSWORD}" tukem-db-prod psql -U "${DB_USER}" -d "${DB_NAME}" < backend/init_db.sql

echo ""
echo "=========================================="
echo "Database setup completed!"
echo "=========================================="

# Verify tables
echo ""
echo "Verifying tables..."
docker compose -f docker-compose.prod.yml exec -e PGPASSWORD="${DB_PASSWORD}" tukem-db-prod psql -U "${DB_USER}" -d "${DB_NAME}" -c "\dt"

echo ""
echo "Done!"

