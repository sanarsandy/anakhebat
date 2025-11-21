# Fix Database Container Error

## Langkah Troubleshooting

### 1. Check Error Logs

```bash
# Check logs database container
docker logs tukem-db-prod

# Atau
docker compose -f docker-compose.prod.yml logs db
```

### 2. Check Container Status

```bash
# Check semua containers
docker ps -a | grep tukem

# Check status
docker compose -f docker-compose.prod.yml ps -a
```

### 3. Stop dan Hapus Semua

```bash
# Stop semua
docker compose -f docker-compose.prod.yml down

# Hapus containers yang stuck
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f 2>/dev/null || true

# Hapus volume (⚠️ Hapus data!)
docker volume ls | grep tukem | awk '{print $2}' | xargs docker volume rm 2>/dev/null || true
```

### 4. Check Port Conflict

```bash
# Check port 5432
sudo netstat -tulpn | grep 5432

# Kill process di port 5432 jika ada
sudo fuser -k 5432/tcp 2>/dev/null || true
```

### 5. Check Environment Variables

```bash
# Pastikan .env file ada
cat .env | grep DB_

# Harus ada:
# DB_HOST=db
# DB_USER=tukem_user
# DB_PASSWORD=your_password
# DB_NAME=tukem_db
# DB_PORT=5432
```

### 6. Start Fresh

```bash
# Start containers
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs -f db
```

### 7. Jika Masih Error, Gunakan PostgreSQL Default

Edit `docker-compose.prod.yml` untuk menggunakan default PostgreSQL:

```yaml
db:
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: tukem_db
```

Dan update `.env`:
```bash
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tukem_db
```

## Quick Fix Script

```bash
#!/bin/bash
# fix-db.sh

echo "Stopping containers..."
docker compose -f docker-compose.prod.yml down

echo "Removing old containers..."
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f 2>/dev/null || true

echo "Killing port 5432..."
sudo fuser -k 5432/tcp 2>/dev/null || true

echo "Starting fresh..."
docker compose -f docker-compose.prod.yml up -d

echo "Waiting for database..."
sleep 10

echo "Checking status..."
docker compose -f docker-compose.prod.yml ps

echo "Database logs:"
docker compose -f docker-compose.prod.yml logs db | tail -20
```

