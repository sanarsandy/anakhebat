# Fix Database Error: "database tukem_user does not exist"

## Problem
Error: `FATAL: database "tukem_user" does not exist`

Ini terjadi karena healthcheck mencoba connect tanpa specify database name yang benar.

## Solution

### Step 1: Pull Latest Code

```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

### Step 2: Check .env File

Pastikan file `.env` sudah benar:

```bash
cat .env
```

Harus ada:
```bash
DB_HOST=db
DB_USER=tukem_user
DB_PASSWORD=your_password
DB_NAME=tukem_db  # <-- Ini harus tukem_db, bukan tukem_user
DB_PORT=5432
```

### Step 3: Restart Containers

```bash
# Stop containers
docker compose -f docker-compose.prod.yml down

# Start ulang
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs db | tail -20
```

### Step 4: Verify Database

```bash
# Check apakah database tukem_db ada
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "\l" | grep tukem

# Atau check langsung
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT version();"
```

### Step 5: Run Migrations

Setelah database ready, jalankan migrations:

```bash
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1 | grep -v "already exists" || true
    echo ""
done
```

## Alternative: Fresh Start

Jika masih ada masalah, fresh start:

```bash
# Stop dan hapus semua
docker compose -f docker-compose.prod.yml down -v

# Hapus volume yang tersisa
docker volume ls | grep tukem | awk '{print $2}' | xargs docker volume rm 2>/dev/null || true

# Start fresh
docker compose -f docker-compose.prod.yml up -d

# Tunggu database ready (10-20 detik)
sleep 15

# Check status
docker compose -f docker-compose.prod.yml ps

# Run migrations
for file in backend/migrations/*.sql; do
    echo "Running: $(basename $file)"
    docker compose -f docker-compose.prod.yml exec -T tukem-db-prod psql -U tukem_user -d tukem_db < "$file" 2>&1 | grep -v "already exists" || true
    echo ""
done
```

