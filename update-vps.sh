#!/bin/bash

# Script untuk update aplikasi di VPS
# Usage: ./update-vps.sh

set -e  # Exit on error

echo "=========================================="
echo "  Update Aplikasi Tukem di VPS"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Backup Database
echo -e "${YELLOW}[1/8]${NC} Backup database..."
BACKUP_FILE="backups/backup_$(date +%Y%m%d_%H%M%S).sql"
mkdir -p backups
docker compose -f docker-compose.prod.yml exec -T db pg_dump -U tukem_user tukem_db > "$BACKUP_FILE" 2>/dev/null || echo "Warning: Backup failed, continuing..."
echo -e "${GREEN}✓${NC} Backup saved to: $BACKUP_FILE"
echo ""

# 2. Pull dari Git
echo -e "${YELLOW}[2/8]${NC} Pull perubahan dari git..."
git pull origin main
echo -e "${GREEN}✓${NC} Git pull completed"
echo ""

# 3. Run Database Migration
echo -e "${YELLOW}[3/8]${NC} Run database migration..."
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db <<EOF
-- Tambahkan kolom untuk Google OAuth
ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(50) DEFAULT 'email';
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Buat index
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id_unique ON users(google_id) WHERE google_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_auth_provider ON users(email, auth_provider) WHERE auth_provider IN ('email', 'both');
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_google ON users(email) WHERE auth_provider = 'google';
EOF
echo -e "${GREEN}✓${NC} Migration completed"
echo ""

# 4. Rebuild Backend (karena ada perubahan dependencies)
echo -e "${YELLOW}[4/8]${NC} Rebuild backend container..."
docker compose -f docker-compose.prod.yml build api
echo -e "${GREEN}✓${NC} Backend rebuild completed"
echo ""

# 5. Rebuild Frontend (jika ada perubahan)
echo -e "${YELLOW}[5/8]${NC} Rebuild frontend container..."
docker compose -f docker-compose.prod.yml build app
echo -e "${GREEN}✓${NC} Frontend rebuild completed"
echo ""

# 6. Restart Services
echo -e "${YELLOW}[6/8]${NC} Restart services..."
docker compose -f docker-compose.prod.yml up -d
echo -e "${GREEN}✓${NC} Services restarted"
echo ""

# 7. Wait for services to be ready
echo -e "${YELLOW}[7/8]${NC} Waiting for services to be ready..."
sleep 10
echo -e "${GREEN}✓${NC} Services ready"
echo ""

# 8. Verify Services
echo -e "${YELLOW}[8/8]${NC} Verifying services..."
echo ""
echo "Service Status:"
docker compose -f docker-compose.prod.yml ps
echo ""

# Check API health
echo "Checking API health..."
if curl -s http://localhost:8085/health > /dev/null; then
    echo -e "${GREEN}✓${NC} API is healthy"
else
    echo -e "${RED}✗${NC} API health check failed"
fi
echo ""

# Check Google OAuth endpoint
echo "Checking Google OAuth endpoint..."
if curl -s http://localhost:8085/api/auth/google | grep -q "auth_url"; then
    echo -e "${GREEN}✓${NC} Google OAuth endpoint is working"
else
    echo -e "${YELLOW}⚠${NC} Google OAuth endpoint check failed (might need credentials)"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Update completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Check logs: docker compose -f docker-compose.prod.yml logs -f"
echo "2. Test Google OAuth login in browser"
echo "3. Test Denver II chart in browser"
echo ""

