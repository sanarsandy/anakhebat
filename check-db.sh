#!/bin/bash

echo "=========================================="
echo "Checking Database Container Status"
echo "=========================================="

echo ""
echo "1. Container Status:"
docker ps -a | grep tukem-db-prod

echo ""
echo "2. Container Logs (last 30 lines):"
docker logs tukem-db-prod 2>&1 | tail -30

echo ""
echo "3. Health Check:"
docker inspect tukem-db-prod --format='{{.State.Health.Status}}' 2>/dev/null || echo "No health check"

echo ""
echo "4. Try to connect:"
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;" 2>&1 || echo "Cannot connect"

echo ""
echo "5. Check if port 5432 is listening:"
docker compose -f docker-compose.prod.yml exec tukem-db-prod netstat -tuln | grep 5432 || echo "Port not listening"

echo ""
echo "=========================================="

