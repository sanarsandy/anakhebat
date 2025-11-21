# Fix Vite Allowed Hosts Error

## Problem
Error saat mengakses `anakhebat.web.id`:
```
Blocked request. This host ("anakhebat.web.id") is not allowed.
To allow this host, add "anakhebat.web.id" to `server.allowedHosts` in vite.config.js.
```

## Solution

### 1. Update nuxt.config.ts

File `frontend/nuxt.config.ts` sudah di-update untuk include:
- `anakhebat.web.id`
- `www.anakhebat.web.id`
- `103.127.134.107`

### 2. Rebuild Frontend Container

Di VPS, jalankan:

```bash
cd /var/rumah_afiat/tukem

# Pull latest code
git pull origin main

# Rebuild frontend container
docker compose -f docker-compose.prod.yml build app

# Restart app container
docker compose -f docker-compose.prod.yml restart app
```

### 3. Alternative: Update via Environment Variable

Jika ingin lebih fleksibel, tambahkan di `.env`:

```env
VITE_ALLOWED_HOSTS=localhost,127.0.0.1,anakhebat.web.id,www.anakhebat.web.id,103.127.134.107
```

Lalu restart container:
```bash
docker compose -f docker-compose.prod.yml restart app
```

### 4. Verify

Setelah rebuild dan restart, akses:
- https://anakhebat.web.id

Error seharusnya sudah hilang.

## Note

Jika masih error, check logs:
```bash
docker compose -f docker-compose.prod.yml logs -f app
```

Pastikan container sudah menggunakan config baru.

