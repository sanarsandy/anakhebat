#!/bin/bash

# Simple database setup script
# Usage: ./setup-database.sh

set -e

echo "=========================================="
echo "Setting up Tukem Database"
echo "=========================================="

# Check if container is running
if ! docker ps | grep -q "tukem-db-prod"; then
    echo "Error: Database container is not running!"
    echo "Please start containers first: docker compose -f docker-compose.prod.yml up -d"
    exit 1
fi

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
echo "Done! Now restart API container to seed data:"
echo "docker compose -f docker-compose.prod.yml restart api"

