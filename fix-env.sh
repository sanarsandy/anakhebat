#!/bin/bash

# Script untuk fix .env file dengan password yang mengandung karakter khusus

echo "=========================================="
echo "Fixing .env File"
echo "=========================================="

if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    exit 1
fi

# Backup .env
cp .env .env.backup
echo "Backup created: .env.backup"

# Read password from .env
PASSWORD=$(grep "^DB_PASSWORD=" .env | cut -d'=' -f2)

if [ -z "$PASSWORD" ]; then
    echo "Error: DB_PASSWORD not found in .env"
    exit 1
fi

# URL encode password (replace / with %2F)
ENCODED_PASSWORD=$(echo "$PASSWORD" | sed 's/\//%2F/g' | sed 's/+/%2B/g' | sed 's/=/%3D/g')

# Update DATABASE_URL
DB_USER=$(grep "^DB_USER=" .env | cut -d'=' -f2)
DB_NAME=$(grep "^DB_NAME=" .env | cut -d'=' -f2)
DB_HOST=$(grep "^DB_HOST=" .env | cut -d'=' -f2)
DB_PORT=$(grep "^DB_PORT=" .env | cut -d'=' -f2)

NEW_DATABASE_URL="postgresql://${DB_USER}:${ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Update .env file
sed -i "s|^DATABASE_URL=.*|DATABASE_URL=${NEW_DATABASE_URL}|" .env

echo ""
echo "✓ DATABASE_URL updated with URL-encoded password"
echo ""
echo "Updated DATABASE_URL:"
grep "^DATABASE_URL=" .env
echo ""
echo "Done!"

