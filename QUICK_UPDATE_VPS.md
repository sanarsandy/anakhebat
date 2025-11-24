# Quick Update di VPS

## Cara Cepat (Menggunakan Script)

```bash
# 1. SSH ke VPS
ssh user@your-vps-ip

# 2. Masuk ke direktori project
cd /var/rumah_afiat/tukem

# 3. Jalankan script update
./update-vps.sh
```

## Cara Manual

```bash
# 1. SSH ke VPS
ssh user@your-vps-ip

# 2. Masuk ke direktori project
cd /var/rumah_afiat/tukem

# 3. Pull dari git
git pull origin main

# 4. Run migration
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user -d tukem_db < backend/migrations/007_add_google_oauth.sql

# 5. Rebuild dan restart
docker compose -f docker-compose.prod.yml build api app
docker compose -f docker-compose.prod.yml up -d

# 6. Cek status
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs api --tail=20
```

## Set Environment Variables (Jika Belum)

Edit `docker-compose.prod.yml` atau buat file `.env`:

```bash
# Di server VPS, edit docker-compose.prod.yml atau set di .env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback
FRONTEND_URL=https://yourdomain.com
JWT_SECRET=your-secure-jwt-secret
```

## Verifikasi Update

```bash
# Test API
curl http://localhost:8085/health

# Test Google OAuth
curl http://localhost:8085/api/auth/google

# Cek logs
docker compose -f docker-compose.prod.yml logs api --tail=50
```

