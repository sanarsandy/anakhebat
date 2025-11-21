# Check dan Fix Missing Tables

## Problem
Migrations sudah berjalan tapi seed data error karena tables tidak ada.

## Solution

### Step 1: Check Tables yang Ada

```bash
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"
```

### Step 2: Run Migrations Manual untuk Memastikan Semua Tables Ada

```bash
cd /var/rumah_afiat/tukem

# Run semua migrations lagi (akan skip yang sudah ada)
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1
    echo ""
done
```

### Step 3: Check Specific Tables

```bash
# Check milestones
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM milestones;" 2>&1

# Check who_standards
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM who_standards;" 2>&1

# Check stimulation_content
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM stimulation_content;" 2>&1

# Check immunization_schedule
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM immunization_schedule;" 2>&1
```

### Step 4: Jika Tables Masih Tidak Ada, Run Migrations Satu Per Satu

```bash
# Migration 004 - Milestones (penting!)
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/004_milestones_tables.sql

# Migration 005 - WHO Standards
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/005_who_standards.sql

# Migration 007 - Stimulation Content
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/007_stimulation_content.sql

# Migration 008 - Immunization (sudah ada, tapi run lagi untuk memastikan)
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/008_immunization_tables.sql
```

### Step 5: Restart API Container untuk Seed Data

Setelah semua tables ada, restart API untuk seed data:

```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs untuk seed data
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

### Step 6: Verify All Tables

```bash
# List semua tables
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"

# Check data di setiap table
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 'users' as table_name, COUNT(*) as count FROM users UNION ALL SELECT 'children', COUNT(*) FROM children UNION ALL SELECT 'measurements', COUNT(*) FROM measurements UNION ALL SELECT 'milestones', COUNT(*) FROM milestones UNION ALL SELECT 'who_standards', COUNT(*) FROM who_standards UNION ALL SELECT 'stimulation_content', COUNT(*) FROM stimulation_content UNION ALL SELECT 'immunization_schedule', COUNT(*) FROM immunization_schedule;"
```

## Expected Tables

Setelah migrations berhasil, harus ada:
- users
- children
- measurements
- milestones
- assessments
- who_standards
- stimulation_content
- immunization_schedule
- child_immunizations

