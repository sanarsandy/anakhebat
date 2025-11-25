#!/bin/bash

# Script untuk menjalankan migrations secara manual
# Usage: ./run-migrations.sh [docker|local]
#   docker: Run migrations via Docker (default if psql not found)
#   local: Run migrations directly using psql

echo "Running database migrations..."

# List migrations in order
MIGRATIONS=(
    "001_init_schema.sql"
    "002_children_table.sql"
    "003_measurements_table.sql"
    "004_milestones_tables.sql"
    "005_who_standards.sql"
    "006_add_denver_domain.sql"
    "007_add_google_oauth.sql"
    "007_stimulation_content.sql"
    "008_immunization_tables.sql"
    "009_add_phone_otp_auth.sql"
)

# Database connection (adjust as needed)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-tukem_user}
DB_NAME=${DB_NAME:-tukem_db}
DB_PASSWORD=${DB_PASSWORD:-tukem_password}

# Docker container name (adjust as needed)
DB_CONTAINER=${DB_CONTAINER:-tukem-db}

# Detect if we should use Docker
USE_DOCKER=false
if [ "$1" = "docker" ]; then
    USE_DOCKER=true
elif [ "$1" = "local" ]; then
    USE_DOCKER=false
else
    # Auto-detect: check if psql is available
    if ! command -v psql &> /dev/null; then
        echo "psql not found, using Docker..."
        USE_DOCKER=true
    else
        USE_DOCKER=false
    fi
fi

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$SCRIPT_DIR/../migrations"

# Run each migration
for migration in "${MIGRATIONS[@]}"; do
    echo "Running migration: $migration"
    
    if [ "$USE_DOCKER" = true ]; then
        # Run via Docker
        if docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$MIGRATIONS_DIR/$migration" 2>&1; then
            echo "✓ Migration $migration completed"
        else
            echo "✗ Migration $migration failed"
            # Check if it's just a "already exists" error (which is OK)
            if docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$MIGRATIONS_DIR/$migration" 2>&1 | grep -q "already exists\|duplicate"; then
                echo "  (Migration already applied, skipping...)"
            fi
        fi
    else
        # Run directly using psql
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$MIGRATIONS_DIR/$migration"
        if [ $? -eq 0 ]; then
            echo "✓ Migration $migration completed"
        else
            echo "✗ Migration $migration failed"
        fi
    fi
done

echo "All migrations completed!"

