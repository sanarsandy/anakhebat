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

# Replace {{PROJECT_NAME}} placeholder in all files
find "$PROJECT_DIR" -type f \( -name "*.yml" -o -name "*.yaml" -o -name "*.go" -o -name "*.ts" -o -name "*.vue" -o -name "*.json" -o -name "*.md" -o -name "*.sql" \) -exec sed -i.bak "s/{{PROJECT_NAME}}/${PROJECT_NAME}/g" {} \;

# Remove backup files
find "$PROJECT_DIR" -type f -name "*.bak" -delete

# Also replace tukem references (legacy)
find "$PROJECT_DIR" -type f \( -name "*.go" -o -name "*.json" \) -exec sed -i.bak "s/tukem-backend/${PROJECT_NAME}-backend/g; s/tukem-frontend/${PROJECT_NAME}-frontend/g" {} \;
find "$PROJECT_DIR" -type f -name "*.bak" -delete

echo "✅ Project generated successfully!"
echo ""
echo "📝 Next steps:"
echo "  1. cd $PROJECT_DIR"
echo "  2. cp .env.example .env"
echo "  3. Edit .env with your configuration"
echo "  4. docker compose up -d"
echo ""
echo "🎉 Happy coding!"

