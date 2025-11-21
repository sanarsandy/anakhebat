#!/bin/bash

echo "=========================================="
echo "Testing Database Connection"
echo "=========================================="

echo ""
echo "1. Check container status:"
docker ps -a | grep tukem-db-prod

echo ""
echo "2. Check if container is running:"
if docker ps | grep -q "tukem-db-prod"; then
    echo "✓ Container is running"
else
    echo "✗ Container is NOT running"
    echo "Start it: docker compose -f docker-compose.prod.yml up -d db"
    exit 1
fi

echo ""
echo "3. Try connect to postgres database:"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;" 2>&1

echo ""
echo "4. Try connect to tukem_db database:"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;" 2>&1

echo ""
echo "5. Check databases:"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "\l" 2>&1 | grep tukem

echo ""
echo "6. Check environment variables:"
docker compose -f docker-compose.prod.yml exec tukem-db-prod env | grep POSTGRES

echo ""
echo "Done!"

