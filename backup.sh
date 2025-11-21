#!/bin/bash

# Backup script untuk database Tukem
# Usage: ./backup.sh

BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="tukem_backup_$DATE.sql"
COMPOSE_FILE="docker-compose.prod.yml"

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Check if docker-compose file exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo "Error: $COMPOSE_FILE not found!"
    exit 1
fi

echo "Starting backup at $(date)"

# Backup database
docker compose -f $COMPOSE_FILE exec -T db pg_dump -U tukem_user tukem_db > $BACKUP_DIR/$FILENAME

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo "Database backup successful: $FILENAME"
    
    # Compress backup
    gzip $BACKUP_DIR/$FILENAME
    echo "Backup compressed: $FILENAME.gz"
    
    # Delete backups older than 30 days
    find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
    echo "Old backups cleaned (older than 30 days)"
    
    echo "Backup completed successfully at $(date)"
else
    echo "Error: Backup failed!"
    exit 1
fi

