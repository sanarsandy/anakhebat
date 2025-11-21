# Deployment Guide - Simple & Clean

## Quick Start

### 1. Setup Environment

```bash
# Copy .env.example to .env
cp .env.example .env

# Edit .env dengan values yang benar
nano .env
```

**Important:** Pastikan password tidak mengandung karakter khusus yang bermasalah seperti `/`, `+`, `=` di DATABASE_URL. Jika ada, gunakan password yang lebih simple.

### 2. Start Database

```bash
# Start database saja dulu
docker compose -f docker-compose.prod.yml up -d db

# Check logs
docker compose -f docker-compose.prod.yml logs -f db
```

### 3. Setup Tables

```bash
# Setup tables
chmod +x setup-db-simple.sh
./setup-db-simple.sh
```

### 4. Start All Services

```bash
# Start semua services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps
```

### 5. Restart API untuk Seed Data

```bash
docker compose -f docker-compose.prod.yml restart api

# Check logs
docker compose -f docker-compose.prod.yml logs -f api
```

## Troubleshooting

### Database tidak bisa connect

```bash
# Test connection
docker exec tukem-db-prod psql -U tukem_user -d tukem_db -c "SELECT 1;"
```

### Password mengandung karakter khusus

Jika password mengandung `/`, `+`, `=`, dll:
1. Generate password baru: `openssl rand -base64 24 | tr -d '/+='`
2. Update di `.env`
3. Restart containers

### Fresh Start

```bash
# Stop dan hapus semua
docker compose -f docker-compose.prod.yml down -v

# Start fresh
docker compose -f docker-compose.prod.yml up -d db
sleep 15
./setup-db-simple.sh
docker compose -f docker-compose.prod.yml up -d
```

## Environment Variables

Lihat `.env.example` untuk format yang benar.

**Recommended Password Format:**
- Gunakan alphanumeric + beberapa karakter khusus yang aman
- Hindari: `/`, `+`, `=`, `@`, `#`, `%`
- Atau gunakan: `openssl rand -base64 24 | tr -d '/+='`

