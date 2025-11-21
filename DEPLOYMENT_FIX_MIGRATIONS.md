# Fix Database Migrations Issue

## Problem
Database tables tidak ada karena migrations tidak dijalankan dengan lengkap.

## Solution

### Option 1: Restart Container (Recommended)

Migrations sekarang sudah diperbaiki di `main.go` untuk menjalankan semua migrations secara berurutan. Cukup restart container:

```bash
# Stop containers
docker compose -f docker-compose.prod.yml down

# Start containers (migrations akan otomatis dijalankan)
docker compose -f docker-compose.prod.yml up -d

# Check logs untuk memastikan migrations berhasil
docker compose -f docker-compose.prod.yml logs api | grep -i migration
```

### Option 2: Run Migrations Manual via Docker

Jika migrations tidak otomatis jalan, jalankan manual:

```bash
# Masuk ke container database
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db

# Atau run migrations via psql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/001_init_schema.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/002_children_table.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/003_measurements_table.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/004_milestones_tables.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/005_who_standards.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/006_add_denver_domain.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/007_stimulation_content.sql
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/008_immunization_tables.sql
```

### Option 3: Run All Migrations dengan Script

```bash
# Copy script ke container atau run dari host
for file in backend/migrations/*.sql; do
    echo "Running: $file"
    docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < "$file"
done
```

### Option 4: Fresh Start (Hapus Semua Data)

⚠️ **WARNING**: Ini akan menghapus semua data!

```bash
# Stop dan hapus semua containers dan volumes
docker compose -f docker-compose.prod.yml down -v

# Start fresh
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs -f api
```

## Verify Migrations

Setelah migrations berjalan, verify dengan:

```bash
# Check tables
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db -c "\dt"

# Check specific tables
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM milestones;"
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM who_standards;"
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM stimulation_content;"
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db -c "SELECT COUNT(*) FROM immunization_schedule;"
```

## Expected Tables

Setelah migrations berhasil, harus ada tables:
- users
- children
- measurements
- milestones
- assessments
- who_standards
- stimulation_content
- immunization_schedule
- child_immunizations

