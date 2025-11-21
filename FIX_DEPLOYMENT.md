# Fix Deployment Issues

## Problem: Container tidak running

### Step 1: Check Status Containers

```bash
cd /var/rumah_afiat/tukem

# Check semua containers
docker compose -f docker-compose.prod.yml ps

# Check semua containers (termasuk yang stopped)
docker ps -a | grep tukem
```

### Step 2: Pull Latest Code (Pastikan docker-compose.prod.yml sudah di-update)

```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

### Step 3: Stop dan Start Ulang Containers

```bash
# Stop semua containers
docker compose -f docker-compose.prod.yml down

# Start ulang
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps
```

### Step 4: Check Logs jika ada masalah

```bash
# Check database logs
docker compose -f docker-compose.prod.yml logs db

# Check API logs
docker compose -f docker-compose.prod.yml logs api

# Check semua logs
docker compose -f docker-compose.prod.yml logs
```

### Step 5: Pastikan Environment Variables Sudah Di-set

```bash
# Check apakah file .env ada
ls -la .env

# Jika belum ada, buat dari template
cat .env
```

### Step 6: Run Migrations Setelah Container Running

Setelah containers running, jalankan migrations:

```bash
# Option 1: Loop semua migrations
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1 | grep -v "already exists" || true
    echo ""
done

# Option 2: Manual satu per satu
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/001_init_schema.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/002_children_table.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/003_measurements_table.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/004_milestones_tables.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/005_who_standards.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/006_add_denver_domain.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/007_stimulation_content.sql
docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < backend/migrations/008_immunization_tables.sql
```

### Step 7: Verifikasi

```bash
# Check tables
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"

# Check container status
docker compose -f docker-compose.prod.yml ps
```

## Troubleshooting

### Jika database container tidak start:

```bash
# Check database logs
docker compose -f docker-compose.prod.yml logs db

# Check apakah port 5432 sudah digunakan
sudo netstat -tulpn | grep 5432

# Check apakah volume ada masalah
docker volume ls | grep tukem
```

### Jika masih ada warning tentang version:

Pastikan sudah pull latest code:
```bash
git pull origin main
cat docker-compose.prod.yml | head -5
# Harusnya tidak ada "version: '3.8'"
```

### Fresh Start (Hapus Semua dan Mulai Lagi)

⚠️ **WARNING**: Ini akan hapus semua data!

```bash
# Stop dan hapus semua
docker compose -f docker-compose.prod.yml down -v

# Hapus containers yang tersisa
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f

# Start fresh
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps
```

