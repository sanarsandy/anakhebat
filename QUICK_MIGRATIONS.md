# Quick Migration Commands untuk VPS

Jika file `run-migrations.sh` tidak ada, jalankan migrations manual dengan commands berikut:

## Step 1: Pull Latest Code

```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

## Step 2: Run Migrations Manual

Jalankan satu per satu:

```bash
# Migration 001 - Init schema (users, children)
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/001_init_schema.sql

# Migration 002 - Children table (bisa skip jika sudah ada dari 001)
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/002_children_table.sql 2>&1 | grep -v "already exists" || true

# Migration 003 - Measurements
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/003_measurements_table.sql

# Migration 004 - Milestones
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/004_milestones_tables.sql

# Migration 005 - WHO Standards
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/005_who_standards.sql

# Migration 006 - Denver Domain
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/006_add_denver_domain.sql

# Migration 007 - Stimulation Content
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/007_stimulation_content.sql

# Migration 008 - Immunization
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/008_immunization_tables.sql
```

## Step 3: Atau Run Semua Sekaligus dengan Loop

```bash
cd /var/rumah_afiat/tukem

for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1 | grep -v "already exists" || true
    echo ""
done
```

## Step 4: Verifikasi

```bash
# Check semua tables
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"

# Check specific tables
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM users;"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM children;"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM milestones;"
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM who_standards;"
```

## Step 5: Restart API untuk Seed Data

```bash
docker compose -f docker-compose.prod.yml restart api
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

