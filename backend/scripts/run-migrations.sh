#!/bin/bash

# Script untuk menjalankan migrations secara manual
# Usage: ./run-migrations.sh

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

# Run each migration
for migration in "${MIGRATIONS[@]}"; do
    echo "Running migration: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "migrations/$migration"
    if [ $? -eq 0 ]; then
        echo "✓ Migration $migration completed"
    else
        echo "✗ Migration $migration failed"
    fi
done

echo "All migrations completed!"

