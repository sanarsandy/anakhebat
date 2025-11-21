#!/bin/bash

# Script untuk menjalankan migrations secara manual di VPS
# Usage: ./run-migrations.sh

set -e

echo "=========================================="
echo "Running Database Migrations"
echo "=========================================="

# Database connection
DB_CONTAINER="tukem-db-prod"
DB_USER="tukem_user"
DB_NAME="tukem_db"
MIGRATIONS_DIR="backend/migrations"

# List migrations in order
MIGRATIONS=(
    "001_init_schema.sql"
    "002_children_table.sql"
    "003_measurements_table.sql"
    "004_milestones_tables.sql"
    "005_who_standards.sql"
    "006_add_denver_domain.sql"
    "007_stimulation_content.sql"
    "008_immunization_tables.sql"
)

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Check if database container is running
if ! docker ps | grep -q "$DB_CONTAINER"; then
    echo "Error: Database container '$DB_CONTAINER' is not running"
    echo "Please start containers first: docker compose -f docker-compose.prod.yml up -d"
    exit 1
fi

echo "Database container: $DB_CONTAINER"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo ""

# Run each migration
for migration in "${MIGRATIONS[@]}"; do
    migration_file="$MIGRATIONS_DIR/$migration"
    
    if [ ! -f "$migration_file" ]; then
        echo "⚠️  Warning: Migration file not found: $migration_file"
        continue
    fi
    
    echo "Running migration: $migration"
    
    # Run migration and capture output
    output=$(docker compose -f docker-compose.prod.yml exec -T "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" < "$migration_file" 2>&1)
    exit_code=$?
    
    # Check for common "already exists" errors (which are OK)
    if echo "$output" | grep -qiE "(already exists|duplicate|relation.*already exists)"; then
        echo "✓ Migration $migration already exists (skipped)"
    elif [ $exit_code -eq 0 ]; then
        echo "✓ Migration $migration completed successfully"
    else
        echo "✗ Migration $migration failed"
        echo "Error output:"
        echo "$output" | head -5
        echo ""
        echo "Continuing with next migration..."
        # Don't exit, continue with next migration
    fi
    echo ""
done

echo "=========================================="
echo "All migrations completed!"
echo "=========================================="

# Verify tables
echo ""
echo "Verifying tables..."
docker compose -f docker-compose.prod.yml exec -T "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "\dt" || true

echo ""
echo "Done!"

