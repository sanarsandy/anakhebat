#!/bin/bash

# Script to generate new project from template
# Usage: ./generate-project.sh <project-name>

set -e

PROJECT_NAME=$1

if [ -z "$PROJECT_NAME" ]; then
    echo "Error: Project name is required"
    echo "Usage: ./generate-project.sh <project-name>"
    exit 1
fi

# Validate project name (alphanumeric and hyphens only)
if [[ ! "$PROJECT_NAME" =~ ^[a-z0-9-]+$ ]]; then
    echo "Error: Project name must contain only lowercase letters, numbers, and hyphens"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")/template"
PROJECT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")/$PROJECT_NAME"

if [ -d "$PROJECT_DIR" ]; then
    echo "Error: Directory $PROJECT_DIR already exists"
    exit 1
fi

echo "🚀 Generating project: $PROJECT_NAME"
echo "📁 Template: $TEMPLATE_DIR"
echo "📁 Destination: $PROJECT_DIR"

# Copy template
echo "📋 Copying template files..."
cp -r "$TEMPLATE_DIR" "$PROJECT_DIR"

# Replace placeholders
echo "🔧 Replacing placeholders..."

# Backend
if [ -f "$PROJECT_DIR/backend/go.mod" ]; then
    sed -i.bak "s/tukem-backend/${PROJECT_NAME}-backend/g" "$PROJECT_DIR/backend/go.mod"
    rm "$PROJECT_DIR/backend/go.mod.bak"
fi

# Frontend
if [ -f "$PROJECT_DIR/frontend/package.json" ]; then
    sed -i.bak "s/tukem-frontend/${PROJECT_NAME}-frontend/g" "$PROJECT_DIR/frontend/package.json"
    rm "$PROJECT_DIR/frontend/package.json.bak"
fi

# Docker Compose
for file in "$PROJECT_DIR"/docker-compose*.yml; do
    if [ -f "$file" ]; then
        sed -i.bak "s/tukem/${PROJECT_NAME}/g" "$file"
        rm "${file}.bak"
    fi
done

echo "✅ Project generated successfully!"
echo ""
echo "📝 Next steps:"
echo "  1. cd $PROJECT_DIR"
echo "  2. cp .env.example .env"
echo "  3. Edit .env with your configuration"
echo "  4. docker compose up -d"
echo ""
echo "🎉 Happy coding!"

