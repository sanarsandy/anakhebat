# Fix Migration Order Issue

## Problem
Migration 001 gagal membuat table `users`, sehingga semua migration berikutnya gagal karena dependency.

## Root Cause
Migration 001 mungkin tidak berjalan dengan benar atau ada error yang tidak terlihat.

## Solution

### Step 1: Check Apakah Table Users Sudah Ada

```bash
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt" | grep users
```

### Step 2: Run Migration 001 Manual (PENTING!)

```bash
cd /var/rumah_afiat/tukem

# Run migration 001 dulu
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/001_init_schema.sql
```

### Step 3: Verify Table Users dan Children Sudah Dibuat

```bash
# Check users table
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM users;"

# Check children table
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM children;"
```

### Step 4: Run Semua Migrations Lagi (Sekarang Harusnya Berhasil)

```bash
# Run semua migrations secara berurutan
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1 | grep -v "already exists" || true
    echo ""
done
```

### Step 5: Verify Semua Tables

```bash
# List semua tables
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"

# Check data di setiap table
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "
SELECT 'users' as table_name, COUNT(*) as count FROM users 
UNION ALL SELECT 'children', COUNT(*) FROM children 
UNION ALL SELECT 'measurements', COUNT(*) FROM measurements 
UNION ALL SELECT 'milestones', COUNT(*) FROM milestones 
UNION ALL SELECT 'assessments', COUNT(*) FROM assessments 
UNION ALL SELECT 'who_standards', COUNT(*) FROM who_standards 
UNION ALL SELECT 'stimulation_content', COUNT(*) FROM stimulation_content 
UNION ALL SELECT 'immunization_schedule', COUNT(*) FROM immunization_schedule 
UNION ALL SELECT 'child_immunizations', COUNT(*) FROM child_immunizations;
"
```

### Step 6: Restart API untuk Seed Data

```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs untuk seed data
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

## Alternative: Fresh Start (Jika Masih Error)

Jika masih ada masalah, fresh start:

```bash
# Stop dan hapus semua
docker compose -f docker-compose.prod.yml down -v

# Hapus volume
docker volume ls | grep tukem | awk '{print $2}' | xargs docker volume rm 2>/dev/null || true

# Start fresh
docker compose -f docker-compose.prod.yml up -d

# Tunggu database ready
sleep 15

# Run migrations
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1
    echo ""
done
```

