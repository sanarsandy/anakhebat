#!/bin/bash

# Script untuk setup tables setelah database running
# Usage: ./setup-tables.sh

set -e

echo "=========================================="
echo "Setting up Database Tables"
echo "=========================================="

# Check if container is running
if ! docker ps | grep -q "tukem-db-prod"; then
    echo "Error: Database container is not running!"
    echo "Please start it first: docker compose -f docker-compose.prod.yml up -d db"
    exit 1
fi

# Wait for database to be ready
echo "Waiting for database to be ready..."
for i in {1..60}; do
    # Try to connect - use postgres database first to avoid connection issues
    if docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    if [ $((i % 5)) -eq 0 ]; then
        echo "Waiting... ($i/60)"
    fi
    sleep 1
done

# Check if database is ready
if ! docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
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
    echo "3. Try manual connection:"
    docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;"
    exit 1
fi

# Run init script
echo ""
echo "Running database initialization..."
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/init_db.sql

echo ""
echo "=========================================="
echo "Database setup completed!"
echo "=========================================="

# Verify tables
echo ""
echo "Verifying tables..."
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"

echo ""
echo "Done!"

