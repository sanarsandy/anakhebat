#!/bin/bash

# Script untuk setup tables setelah database running
# Usage: ./setup-tables.sh

set -e

echo "=========================================="
echo "Setting up Database Tables"
echo "=========================================="

# Wait for database to be ready
echo "Waiting for database to be ready..."
for i in {1..30}; do
    if docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting... ($i/30)"
    sleep 1
done

# Check if database is ready
if ! docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;" > /dev/null 2>&1; then
    echo "Error: Database is not ready after 30 seconds"
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

