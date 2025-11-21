# Final Fix: Database Auto-Setup

## Solusi: Database Setup Otomatis

Sekarang `init_db.sql` akan **otomatis dijalankan** saat database container start pertama kali!

## Langkah Setup

### 1. Pull Latest Code
```bash
cd /var/rumah_afiat/tukem
git pull origin main
```

### 2. Hapus Volume Lama (Fresh Start)
```bash
# Stop semua
docker compose -f docker-compose.prod.yml down

# Hapus volume (⚠️ Hapus data lama!)
docker volume rm tukem_postgres_data 2>/dev/null || true
```

### 3. Start Database (Tables akan dibuat otomatis!)
```bash
# Start database saja dulu
docker compose -f docker-compose.prod.yml up -d db

# Check logs - harusnya lihat tables dibuat
docker compose -f docker-compose.prod.yml logs -f db
```

### 4. Verify Tables
```bash
# Tunggu beberapa detik, lalu check
sleep 10
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt"
```

### 5. Start Semua Services
```bash
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps
```

### 6. Restart API untuk Seed Data
```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs
docker compose -f docker-compose.prod.yml logs -f api | grep -i seed
```

## Cara Kerja

1. PostgreSQL container akan mount `backend/init_db.sql` ke `/docker-entrypoint-initdb.d/01-init.sql`
2. Saat database **pertama kali** start, PostgreSQL akan otomatis menjalankan semua `.sql` files di folder tersebut
3. Tables akan dibuat otomatis!
4. Tidak perlu manual migration lagi!

## Jika Volume Sudah Ada

Jika volume sudah ada, init script tidak akan jalan lagi. Solusinya:

```bash
# Hapus volume
docker compose -f docker-compose.prod.yml down -v

# Start fresh
docker compose -f docker-compose.prod.yml up -d db
```

## Quick Command

```bash
cd /var/rumah_afiat/tukem && \
git pull origin main && \
docker compose -f docker-compose.prod.yml down -v && \
docker compose -f docker-compose.prod.yml up -d db && \
sleep 15 && \
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "\dt" && \
docker compose -f docker-compose.prod.yml up -d
```

