#!/bin/bash

# Simple database setup script
# Usage: ./setup-db-simple.sh

set -e

echo "=========================================="
echo "Simple Database Setup"
echo "=========================================="

# Load .env if exists
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Set defaults
DB_USER=${DB_USER:-tukem_user}
DB_PASSWORD=${DB_PASSWORD:-tukem_password}
DB_NAME=${DB_NAME:-tukem_db}

echo "Database User: $DB_USER"
echo "Database Name: $DB_NAME"
echo ""

# Check if container is running
if ! docker ps | grep -q "tukem-db-prod"; then
    echo "Error: Database container is not running!"
    echo "Start it: docker compose -f docker-compose.prod.yml up -d db"
    exit 1
fi

# Wait for database
echo "Waiting for database..."
for i in {1..30}; do
    if docker exec tukem-db-prod pg_isready -U "$DB_USER" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    sleep 1
done

# Setup tables using PGPASSWORD
echo ""
echo "Setting up tables..."
export PGPASSWORD="$DB_PASSWORD"
docker exec -e PGPASSWORD="$DB_PASSWORD" tukem-db-prod psql -U "$DB_USER" -d "$DB_NAME" -f - < backend/init_db.sql

echo ""
echo "=========================================="
echo "Verifying tables..."
echo "=========================================="
docker exec -e PGPASSWORD="$DB_PASSWORD" tukem-db-prod psql -U "$DB_USER" -d "$DB_NAME" -c "\dt"

echo ""
echo "Done!"

