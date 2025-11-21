# Setup Domain anakhebat.web.id dengan Nginx Reverse Proxy

## Informasi Penting
- **Domain**: anakhebat.web.id
- **IP VPS**: 103.127.134.107
- **Frontend Port**: 3000 (localhost)
- **Backend API Port**: 8085 (localhost)

---

## STEP 1: Setup DNS Records

### 1.1 Login ke Management DNS Provider

Masuk ke panel management DNS provider Anda (misalnya: Cloudflare, Namecheap, GoDaddy, cPanel, dll).

### 1.2 Tambahkan A Records

Tambahkan **2 A Records** berikut:

#### Record 1: Root Domain
```
Type: A
Name: @ (atau kosong, atau anakhebat.web.id)
Value: 103.127.134.107
TTL: 3600 (atau Auto/Default)
Priority: (kosongkan jika tidak ada)
```

#### Record 2: WWW Subdomain
```
Type: A
Name: www
Value: 103.127.134.107
TTL: 3600 (atau Auto/Default)
Priority: (kosongkan jika tidak ada)
```

### 1.3 Contoh Setup di Beberapa Provider

#### Cloudflare
1. Login ke https://dash.cloudflare.com
2. Pilih domain `anakhebat.web.id`
3. Klik tab **"DNS"**
4. Klik **"Add record"**
5. Isi form:
   - **Type**: `A`
   - **Name**: `@`
   - **IPv4 address**: `103.127.134.107`
   - **Proxy status**: `DNS only` (orange cloud OFF)
   - **TTL**: `Auto`
6. Klik **"Save"**
7. Ulangi untuk `www` (Name: `www`, Value: `103.127.134.107`)

#### Namecheap
1. Login ke https://www.namecheap.com
2. **Domain List** → Pilih `anakhebat.web.id` → **Manage** → **Advanced DNS**
3. Di bagian **"Host Records"**, klik **"Add New Record"**
4. Isi:
   - **Type**: `A Record`
   - **Host**: `@`
   - **Value**: `103.127.134.107`
   - **TTL**: `Automatic`
5. Klik **save** (centang)
6. Ulangi untuk `www` (Host: `www`)

#### GoDaddy
1. Login ke https://www.godaddy.com
2. **My Products** → Klik **DNS** di domain `anakhebat.web.id`
3. Klik **"Add"** di bagian Records
4. Isi:
   - **Type**: `A`
   - **Name**: `@`
   - **Value**: `103.127.134.107`
   - **TTL**: `600` (10 minutes) atau `3600` (1 hour)
5. Klik **"Save"**
6. Ulangi untuk `www` (Name: `www`)

#### cPanel
1. Login ke cPanel
2. **Zone Editor** → **Manage**
3. Klik **"Add Record"**
4. Isi:
   - **Type**: `A`
   - **Name**: `@` atau `anakhebat.web.id`
   - **Address**: `103.127.134.107`
   - **TTL**: `14400`
5. Klik **"Add Record"**
6. Ulangi untuk `www` (Name: `www`)

### 1.4 Verifikasi DNS Setup

Tunggu **5-15 menit** (bisa sampai 24 jam untuk full propagation), lalu verifikasi:

#### Dari Terminal VPS
```bash
# Cek A record
dig anakhebat.web.id +short
# Output harus: 103.127.134.107

# Atau pakai nslookup
nslookup anakhebat.web.id
# Output harus menunjukkan: 103.127.134.107

# Atau pakai host
host anakhebat.web.id
```

#### Dari Online Tools
- https://dnschecker.org/#A/anakhebat.web.id
- https://www.whatsmydns.net/#A/anakhebat.web.id
- https://mxtoolbox.com/DNSLookup.aspx

Masukkan domain `anakhebat.web.id` dan pastikan semua server menunjukkan IP `103.127.134.107`.

---

## STEP 2: Install Nginx di VPS

### 2.1 Update Package List
```bash
sudo apt update
```

### 2.2 Install Nginx
```bash
sudo apt install nginx -y
```

### 2.3 Check Status Nginx
```bash
sudo systemctl status nginx
```

Jika status **active (running)**, berarti Nginx sudah terinstall dengan benar.

---

## STEP 3: Buat Nginx Configuration File

### 3.1 Buat File Config
```bash
sudo nano /etc/nginx/sites-available/anakhebat
```

### 3.2 Copy-Paste Konfigurasi Berikut

```nginx
# HTTP Server - Akan diupdate otomatis oleh certbot untuk redirect ke HTTPS
server {
    listen 80;
    server_name anakhebat.web.id www.anakhebat.web.id;

    # Redirect www to non-www
    if ($host = www.anakhebat.web.id) {
        return 301 http://anakhebat.web.id$request_uri;
    }

    # Frontend (Nuxt.js)
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Backend API
    location /api {
        proxy_pass http://127.0.0.1:8085;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:8085/api/health;
        access_log off;
    }

    # Logging
    access_log /var/log/nginx/anakhebat-access.log;
    error_log /var/log/nginx/anakhebat-error.log;
}
```

### 3.3 Save File
- Tekan `Ctrl + O` untuk save
- Tekan `Enter` untuk confirm
- Tekan `Ctrl + X` untuk exit

---

## STEP 4: Enable Nginx Site

### 4.1 Create Symlink
```bash
sudo ln -s /etc/nginx/sites-available/anakhebat /etc/nginx/sites-enabled/
```

### 4.2 Remove Default Site (Optional)
```bash
sudo rm /etc/nginx/sites-enabled/default
```

### 4.3 Test Nginx Configuration
```bash
sudo nginx -t
```

Output harus menunjukkan:
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

Jika ada error, perbaiki file config di Step 3.

### 4.4 Reload Nginx
```bash
sudo systemctl reload nginx
```

---

## STEP 5: Setup SSL dengan Let's Encrypt (HTTPS)

### 5.1 Install Certbot
```bash
sudo apt install certbot python3-certbot-nginx -y
```

### 5.2 Generate SSL Certificate
```bash
sudo certbot --nginx -d anakhebat.web.id -d www.anakhebat.web.id
```

Certbot akan menanyakan beberapa pertanyaan:
1. **Email address**: Masukkan email Anda (untuk notifikasi renewal)
2. **Terms of Service**: Ketik `A` untuk Agree
3. **Share email**: Ketik `Y` untuk Yes atau `N` untuk No
4. **Redirect HTTP to HTTPS**: Ketik `2` untuk Redirect (recommended)

Certbot akan:
- Generate SSL certificate
- Otomatis update Nginx config untuk HTTPS
- Setup auto-renewal

### 5.3 Test Auto-renewal
```bash
sudo certbot renew --dry-run
```

Jika berhasil, berarti SSL akan auto-renew setiap 90 hari.

---

## STEP 6: Configure Firewall

### 6.1 Allow HTTP dan HTTPS
```bash
sudo ufw allow 'Nginx Full'
```

Atau secara manual:
```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### 6.2 Check Firewall Status
```bash
sudo ufw status
```

Pastikan port 80 dan 443 sudah **ALLOW**.

---

## STEP 7: Update Environment Variables

### 7.1 Edit .env File
```bash
cd /var/rumah_afiat/tukem
nano .env
```

### 7.2 Update Konfigurasi Berikut

Cari dan update (atau tambahkan jika belum ada):

```env
# Frontend Configuration
FRONTEND_URL=https://anakhebat.web.id
NUXT_PUBLIC_API_BASE=https://anakhebat.web.id/api

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://anakhebat.web.id,http://anakhebat.web.id
```

### 7.3 Save File
- Tekan `Ctrl + O` untuk save
- Tekan `Enter` untuk confirm
- Tekan `Ctrl + X` untuk exit

---

## STEP 8: Restart Docker Containers

### 8.1 Restart API dan App Containers
```bash
cd /var/rumah_afiat/tukem
docker compose -f docker-compose.prod.yml restart api app
```

### 8.2 Check Container Status
```bash
docker compose -f docker-compose.prod.yml ps
```

Pastikan semua container status **Up**.

### 8.3 Check Logs (Optional)
```bash
# Check API logs
docker compose -f docker-compose.prod.yml logs -f api

# Check App logs
docker compose -f docker-compose.prod.yml logs -f app
```

---

## STEP 9: Verifikasi Setup

### 9.1 Test dari Browser

Buka browser dan akses:
- **Frontend**: https://anakhebat.web.id
- **API Health**: https://anakhebat.web.id/api/health

### 9.2 Test dari Terminal

```bash
# Test frontend
curl -I https://anakhebat.web.id

# Test API health
curl https://anakhebat.web.id/api/health
```

### 9.3 Check Nginx Logs

```bash
# Access log
sudo tail -f /var/log/nginx/anakhebat-access.log

# Error log
sudo tail -f /var/log/nginx/anakhebat-error.log
```

---

## Troubleshooting

### Problem 1: DNS belum resolve

**Gejala**: Domain tidak bisa diakses, atau resolve ke IP yang salah.

**Solusi**:
1. Tunggu lebih lama (bisa sampai 24 jam untuk full propagation)
2. Clear DNS cache di komputer Anda:
   - **Windows**: `ipconfig /flushdns`
   - **Mac**: `sudo dscacheutil -flushcache`
   - **Linux**: `sudo systemd-resolve --flush-caches`
3. Cek di https://dnschecker.org/#A/anakhebat.web.id
4. Pastikan A records sudah benar di DNS management

### Problem 2: Nginx error

**Gejala**: `sudo nginx -t` menunjukkan error.

**Solusi**:
```bash
# Check error detail
sudo nginx -t

# Check Nginx logs
sudo tail -f /var/log/nginx/error.log

# Check config file
sudo nano /etc/nginx/sites-available/anakhebat
# Pastikan syntax benar, tidak ada typo
```

### Problem 3: Port 80 atau 443 sudah digunakan

**Gejala**: Nginx tidak bisa start, atau error "address already in use".

**Solusi**:
```bash
# Check port 80
sudo netstat -tulpn | grep :80

# Check port 443
sudo netstat -tulpn | grep :443

# Kill process jika perlu (ganti <PID> dengan process ID)
sudo kill -9 <PID>

# Atau stop service yang menggunakan port tersebut
sudo systemctl stop apache2  # jika menggunakan Apache
```

### Problem 4: CORS Error

**Gejala**: Browser console menunjukkan CORS error saat akses API.

**Solusi**:
1. Pastikan `CORS_ALLOWED_ORIGINS` di `.env` sudah benar:
   ```env
   CORS_ALLOWED_ORIGINS=https://anakhebat.web.id,http://anakhebat.web.id
   ```
2. Restart API container:
   ```bash
   docker compose -f docker-compose.prod.yml restart api
   ```
3. Check browser console untuk error detail
4. Pastikan frontend menggunakan URL yang sama dengan yang ada di CORS

### Problem 5: SSL Certificate Error

**Gejala**: Browser menunjukkan "Not Secure" atau SSL error.

**Solusi**:
```bash
# Check certificate
sudo certbot certificates

# Renew certificate manual
sudo certbot renew

# Check Nginx config untuk SSL
sudo nginx -t
sudo systemctl reload nginx
```

### Problem 6: 502 Bad Gateway

**Gejala**: Browser menunjukkan "502 Bad Gateway".

**Solusi**:
1. Check apakah containers running:
   ```bash
   docker compose -f docker-compose.prod.yml ps
   ```
2. Check container logs:
   ```bash
   docker compose -f docker-compose.prod.yml logs api
   docker compose -f docker-compose.prod.yml logs app
   ```
3. Restart containers:
   ```bash
   docker compose -f docker-compose.prod.yml restart api app
   ```
4. Check Nginx error log:
   ```bash
   sudo tail -f /var/log/nginx/anakhebat-error.log
   ```

### Problem 7: Frontend tidak bisa connect ke API

**Gejala**: Frontend load tapi API calls gagal.

**Solusi**:
1. Pastikan `NUXT_PUBLIC_API_BASE` di `.env` sudah benar:
   ```env
   NUXT_PUBLIC_API_BASE=https://anakhebat.web.id/api
   ```
2. Restart app container:
   ```bash
   docker compose -f docker-compose.prod.yml restart app
   ```
3. Check browser Network tab untuk melihat request URL
4. Test API langsung:
   ```bash
   curl https://anakhebat.web.id/api/health
   ```

---

## Checklist Final

Setelah semua step selesai, pastikan:

- [ ] DNS A records sudah ditambahkan (@ dan www)
- [ ] DNS sudah resolve ke IP 103.127.134.107
- [ ] Nginx sudah terinstall dan running
- [ ] Nginx config file sudah dibuat dan enabled
- [ ] SSL certificate sudah di-generate dengan certbot
- [ ] Firewall sudah allow port 80 dan 443
- [ ] `.env` file sudah di-update dengan domain baru
- [ ] Docker containers sudah di-restart
- [ ] Frontend bisa diakses di https://anakhebat.web.id
- [ ] API bisa diakses di https://anakhebat.web.id/api/health
- [ ] Tidak ada error di browser console
- [ ] Tidak ada error di Nginx logs

---

## Catatan Penting

1. **DNS Propagation**: Perubahan DNS tidak instant, bisa memakan waktu 5 menit - 24 jam. Gunakan DNS checker online untuk monitor.

2. **SSL Auto-renewal**: Certbot sudah setup auto-renewal. Certificate akan otomatis diperpanjang setiap 90 hari. Pastikan cron job certbot aktif:
   ```bash
   sudo systemctl status certbot.timer
   ```

3. **Backup**: Sebelum melakukan perubahan besar, backup config Nginx:
   ```bash
   sudo cp /etc/nginx/sites-available/anakhebat /etc/nginx/sites-available/anakhebat.backup
   ```

4. **Monitoring**: Setup monitoring untuk Nginx dan containers agar tahu jika ada masalah:
   ```bash
   # Check Nginx status
   sudo systemctl status nginx
   
   # Check containers
   docker compose -f docker-compose.prod.yml ps
   ```

---

## Quick Reference Commands

```bash
# Nginx
sudo systemctl status nginx
sudo systemctl reload nginx
sudo systemctl restart nginx
sudo nginx -t

# SSL
sudo certbot certificates
sudo certbot renew
sudo certbot renew --dry-run

# Docker
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml restart api app
docker compose -f docker-compose.prod.yml logs -f api

# DNS
dig anakhebat.web.id +short
nslookup anakhebat.web.id

# Firewall
sudo ufw status
sudo ufw allow 'Nginx Full'

# Logs
sudo tail -f /var/log/nginx/anakhebat-access.log
sudo tail -f /var/log/nginx/anakhebat-error.log
```

---

**Selesai!** Aplikasi Anda sekarang bisa diakses di **https://anakhebat.web.id** 🎉

