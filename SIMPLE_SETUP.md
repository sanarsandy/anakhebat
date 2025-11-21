# Simple Database Setup

## Database sudah running! Sekarang setup tables:

### Step 1: Pull Latest Code
```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

### Step 2: Setup Database (Buat semua tables)
```bash
# Option 1: Pakai script
./setup-database.sh

# Option 2: Manual
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/init_db.sql
```

### Step 3: Verify Tables
```bash
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"
```

### Step 4: Restart API untuk Seed Data
```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

## Done! ✅

