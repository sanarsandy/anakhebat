# Troubleshoot Container Not Running

## Problem
Container tidak running meskipun sudah di-start.

## Solution

### Step 1: Check Container Status

```bash
# Check semua containers (termasuk yang stopped)
docker ps -a | grep tukem

# Check status dengan docker compose
docker compose -f docker-compose.prod.yml ps
```

### Step 2: Check Logs untuk Error

```bash
# Check database logs
docker compose -f docker-compose.prod.yml logs db

# Check API logs
docker compose -f docker-compose.prod.yml logs api

# Check semua logs
docker compose -f docker-compose.prod.yml logs
```

### Step 3: Check Container yang Stopped

```bash
# List semua containers tukem
docker ps -a | grep tukem

# Check logs container yang stopped
docker logs tukem-db-prod
docker logs tukem-api-prod
docker logs tukem-app-prod
```

### Step 4: Start Container Manual

```bash
# Start database container manual
docker start tukem-db-prod

# Check status
docker ps | grep tukem
```

### Step 5: Check Environment Variables

```bash
# Check apakah .env file ada
ls -la .env

# Check isi .env
cat .env
```

### Step 6: Fresh Start (Jika Masih Error)

```bash
# Stop dan remove semua
docker compose -f docker-compose.prod.yml down

# Remove containers yang tersisa
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f 2>/dev/null || true

# Start fresh
docker compose -f docker-compose.prod.yml up -d

# Check status setelah beberapa detik
sleep 5
docker compose -f docker-compose.prod.yml ps
```

### Step 7: Check Docker Resources

```bash
# Check disk space
df -h

# Check Docker system info
docker system df

# Check jika ada masalah dengan Docker
docker info
```

### Step 8: Check Port Conflicts

```bash
# Check port 5432 (PostgreSQL)
sudo netstat -tulpn | grep 5432

# Check port 8085 (API)
sudo netstat -tulpn | grep 8085

# Check port 3000 (Frontend)
sudo netstat -tulpn | grep 3000
```

## Common Issues

### Issue 1: Container Exit Immediately

Jika container exit segera setelah start, check logs:
```bash
docker logs tukem-db-prod
```

Kemungkinan:
- Environment variables tidak set
- Port conflict
- Volume permission issue

### Issue 2: Healthcheck Failed

Jika healthcheck failed:
```bash
# Check database logs
docker compose -f docker-compose.prod.yml logs db | tail -50

# Test connection manual
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;"
```

### Issue 3: Volume Permission

Jika ada permission issue:
```bash
# Check volume
docker volume ls | grep tukem

# Inspect volume
docker volume inspect tukem_postgres_data
```

