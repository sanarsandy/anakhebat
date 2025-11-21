# Simple Database Setup - Final Solution

## Problem
Init script diabaikan karena volume sudah ada.

## Solution: Setup Tables Setelah Database Running

### Step 1: Pastikan Database Running

```bash
cd /var/rumah_afiat/tukem

# Check database status
docker compose -f docker-compose.prod.yml ps db

# Check logs
docker compose -f docker-compose.prod.yml logs db | tail -10
```

### Step 2: Setup Tables (Setelah Database Running)

```bash
# Pull latest code
git pull origin main

# Run setup script
chmod +x setup-tables.sh
./setup-tables.sh
```

### Atau Manual:

```bash
# Wait for database (check if ready)
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;"

# Run init script
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/init_db.sql

# Verify
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"
```

### Step 3: Restart API untuk Seed Data

```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

## Quick Command

```bash
cd /var/rumah_afiat/tukem && \
git pull origin main && \
./setup-tables.sh && \
docker compose -f docker-compose.prod.yml restart api
```

## Done! ✅

