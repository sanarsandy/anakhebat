# Quick Setup - Database Sudah Running!

## Database sudah running! Sekarang setup tables:

### Step 1: Pull Latest Code
```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

### Step 2: Setup Tables (PENTING!)

```bash
# Option 1: Pakai script
chmod +x setup-db-simple.sh
./setup-db-simple.sh

# Option 2: Manual (jika script error)
export PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q'
docker exec -e PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q' tukem-db-prod psql -U tukem_user -d tukem_db -f - < backend/init_db.sql
```

### Step 3: Verify Tables

```bash
export PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q'
docker exec -e PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q' tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"
```

### Step 4: Restart API untuk Seed Data

```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

## Quick Command (Copy Paste)

```bash
cd /var/rumah_afiat/tukem && \
git pull origin main && \
export PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q' && \
docker exec -e PGPASSWORD="$PGPASSWORD" tukem-db-prod psql -U tukem_user -d tukem_db -f - < backend/init_db.sql && \
docker exec -e PGPASSWORD="$PGPASSWORD" tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt" && \
docker compose -f docker-compose.prod.yml restart api
```

