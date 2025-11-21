# Fix Environment Password Issue

## Problem
Password mengandung karakter khusus `/` yang menyebabkan masalah di:
1. DATABASE_URL (perlu URL encoding)
2. Command line psql
3. Docker environment variables

## Solution

### Option 1: URL Encode Password di DATABASE_URL

Password: `HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q`

Karakter `/` harus di-encode menjadi `%2F`

DATABASE_URL yang benar:
```
DATABASE_URL=postgresql://tukem_user:HySZkeXn9Yx7io4uBvGz9kxyBX0woL%2Fq@db:5432/tukem_db?sslmode=disable
```

### Option 2: Generate Password Baru (Recommended)

```bash
# Generate password baru tanpa karakter khusus
openssl rand -base64 32 | tr -d '/+=' | head -c 32
```

### Option 3: Update .env File

Update file `.env` dengan password yang di-encode:

```bash
# Database Configuration
DB_HOST=db
DB_USER=tukem_user
DB_PASSWORD=HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q
DB_NAME=tukem_db
DB_PORT=5432
DATABASE_URL=postgresql://tukem_user:HySZkeXn9Yx7io4uBvGz9kxyBX0woL%2Fq@db:5432/tukem_db?sslmode=disable
```

## Test Connection

Setelah fix, test connection:

```bash
# Test dengan password yang di-quote
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;"
```

## Alternative: Use PGPASSWORD Environment Variable

```bash
export PGPASSWORD='HySZkeXn9Yx7io4uBvGz9kxyBX0woL/q'
docker compose -f docker-compose.prod.yml exec tukem-db-prod psql -U tukem_user -d postgres -c "SELECT 1;"
```

