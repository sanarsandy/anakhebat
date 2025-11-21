# Panduan Deployment Tukem ke VPS

Panduan lengkap untuk deploy aplikasi Tukem ke VPS menggunakan Docker.

---

## 📋 Prerequisites

### VPS Requirements
- **OS**: Ubuntu 20.04+ / Debian 11+ / CentOS 8+
- **RAM**: Minimum 2GB (recommended 4GB+)
- **Storage**: Minimum 20GB
- **Docker**: Versi 20.10+
- **Docker Compose**: Versi 2.0+

### Software yang Diperlukan
- Docker & Docker Compose (sudah terinstall)
- Git
- Nginx (optional, untuk reverse proxy)
- Certbot (optional, untuk SSL)

---

## 🚀 Step-by-Step Deployment

### Step 1: Persiapan VPS

#### 1.1. Update System
```bash
sudo apt update && sudo apt upgrade -y
```

#### 1.2. Install Git (jika belum ada)
```bash
sudo apt install git -y
```

#### 1.3. Install Docker Compose (jika belum ada)
```bash
# Check Docker version
docker --version

# Check Docker Compose version
docker compose version

# Jika belum ada, install Docker Compose
sudo apt install docker-compose-plugin -y
```

---

### Step 2: Clone Repository

```bash
# Buat directory untuk aplikasi
mkdir -p ~/apps
cd ~/apps

# Clone repository
git clone https://github.com/sanarsandy/anakhebat.git tukem
cd tukem
```

---

### Step 3: Setup Environment Variables

#### 3.1. Buat File `.env` untuk Production

```bash
# Copy dari template (jika ada)
# Atau buat manual
nano .env
```

**Isi file `.env`:**

```bash
# ============================================
# Database Configuration
# ============================================
DB_HOST=db
DB_USER=tukem_user
DB_PASSWORD=CHANGE_THIS_STRONG_PASSWORD
DB_NAME=tukem_db
DB_PORT=5432
DATABASE_URL=postgresql://tukem_user:CHANGE_THIS_STRONG_PASSWORD@db:5432/tukem_db?sslmode=disable

# ============================================
# Server Configuration
# ============================================
PORT=8080
# Note: Port mapping di docker-compose.prod.yml menggunakan 8085:8080
ENV=production

# ============================================
# JWT Configuration
# ============================================
JWT_SECRET=CHANGE_THIS_TO_RANDOM_SECRET_KEY_MIN_32_CHARS
JWT_EXPIRY=72h

# ============================================
# Frontend Configuration
# ============================================
FRONTEND_URL=https://yourdomain.com
NUXT_PUBLIC_API_BASE=https://yourdomain.com/api

# ============================================
# CORS Configuration
# ============================================
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

**⚠️ IMPORTANT:**
- Ganti `CHANGE_THIS_STRONG_PASSWORD` dengan password yang kuat untuk database
- Ganti `CHANGE_THIS_TO_RANDOM_SECRET_KEY_MIN_32_CHARS` dengan secret key random (minimal 32 karakter)
- Ganti `yourdomain.com` dengan domain Anda (atau IP VPS jika tidak pakai domain)

**Generate Random Secret:**
```bash
# Generate random secret untuk JWT
openssl rand -base64 32

# Generate random password untuk database
openssl rand -base64 24
```

---

### Step 4: Update Docker Compose untuk Production

Buat file `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: tukem-db-prod
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - tukem-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: tukem-api-prod
    restart: unless-stopped
    ports:
      - "127.0.0.1:8085:8080"  # Bind ke localhost port 8085, akan di-proxy oleh Nginx
    environment:
      DB_HOST: ${DB_HOST}
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_PORT: ${DB_PORT}
      PORT: ${PORT}
      JWT_SECRET: ${JWT_SECRET}
      ENV: ${ENV}
    depends_on:
      db:
        condition: service_healthy
    networks:
      - tukem-network
    volumes:
      - ./backend:/app
    command: go run main.go

  app:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: tukem-app-prod
    restart: unless-stopped
    ports:
      - "127.0.0.1:3000:3000"  # Bind ke localhost saja, akan di-proxy oleh Nginx
    environment:
      NUXT_PUBLIC_API_BASE: ${NUXT_PUBLIC_API_BASE}
    depends_on:
      - api
    networks:
      - tukem-network
    volumes:
      - ./frontend:/app
      - /app/node_modules

networks:
  tukem-network:
    driver: bridge

volumes:
  postgres_data:
    driver: local
```

---

### Step 5: Build dan Start Services

```bash
# Build images
docker compose -f docker-compose.prod.yml build

# Start services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps

# Check logs
docker compose -f docker-compose.prod.yml logs -f
```

---

### Step 6: Setup Nginx Reverse Proxy (Recommended)

#### 6.1. Install Nginx

```bash
sudo apt install nginx -y
```

#### 6.2. Buat Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/tukem
```

**Isi configuration:**

```nginx
# HTTP to HTTPS redirect
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    # Redirect semua HTTP ke HTTPS
    return 301 https://$server_name$request_uri;
}

# HTTPS Configuration
server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL Certificate (akan di-setup dengan Certbot)
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Frontend (Nuxt App)
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
    }

    # Backend API
    location /api {
        proxy_pass http://127.0.0.1:8085;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers (jika diperlukan)
        add_header Access-Control-Allow-Origin $http_origin always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Authorization, Content-Type" always;
        
        if ($request_method = OPTIONS) {
            return 204;
        }
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:8085/health;
        access_log off;
    }
}
```

**Ganti `yourdomain.com` dengan domain Anda!**

#### 6.3. Enable Site

```bash
# Create symlink
sudo ln -s /etc/nginx/sites-available/tukem /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

### Step 7: Setup SSL Certificate dengan Let's Encrypt

#### 7.1. Install Certbot

```bash
sudo apt install certbot python3-certbot-nginx -y
```

#### 7.2. Get SSL Certificate

```bash
# Ganti yourdomain.com dengan domain Anda
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

Certbot akan:
- Otomatis mengupdate Nginx configuration
- Setup auto-renewal
- Redirect HTTP ke HTTPS

#### 7.3. Test Auto-Renewal

```bash
sudo certbot renew --dry-run
```

---

### Step 8: Firewall Configuration

```bash
# Install UFW (jika belum ada)
sudo apt install ufw -y

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP & HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

---

### Step 9: Verifikasi Deployment

#### 9.1. Check Services

```bash
# Check Docker containers
docker compose -f docker-compose.prod.yml ps

# Check logs
docker compose -f docker-compose.prod.yml logs api
docker compose -f docker-compose.prod.yml logs app
docker compose -f docker-compose.prod.yml logs db
```

#### 9.2. Test Endpoints

```bash
# Health check
curl http://localhost:8085/health

# Frontend
curl http://localhost:3000

# Via domain (jika sudah setup)
curl https://yourdomain.com/health
```

#### 9.3. Test di Browser

- Buka: `https://yourdomain.com`
- Test login/register
- Test semua fitur

---

## 🔧 Maintenance & Management

### Start/Stop Services

```bash
# Start
docker compose -f docker-compose.prod.yml up -d

# Stop
docker compose -f docker-compose.prod.yml down

# Restart
docker compose -f docker-compose.prod.yml restart

# Stop dan hapus volumes (⚠️ HATI-HATI: akan hapus data!)
docker compose -f docker-compose.prod.yml down -v
```

### View Logs

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f api
docker compose -f docker-compose.prod.yml logs -f app
docker compose -f docker-compose.prod.yml logs -f db
```

### Update Application

```bash
# Pull latest code
cd ~/apps/tukem
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# Check logs
docker compose -f docker-compose.prod.yml logs -f
```

### Database Backup

```bash
# Manual backup
docker compose -f docker-compose.prod.yml exec db pg_dump -U tukem_user tukem_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore backup
docker compose -f docker-compose.prod.yml exec -T db psql -U tukem_user tukem_db < backup_20241121_120000.sql
```

### Automated Backup Script

Buat file `backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/home/youruser/apps/tukem/backups"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="tukem_backup_$DATE.sql"

mkdir -p $BACKUP_DIR

docker compose -f docker-compose.prod.yml exec -T db pg_dump -U tukem_user tukem_db > $BACKUP_DIR/$FILENAME

# Compress
gzip $BACKUP_DIR/$FILENAME

# Delete backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: $FILENAME.gz"
```

Setup cron untuk auto-backup:

```bash
# Edit crontab
crontab -e

# Add line untuk backup setiap hari jam 2 pagi
0 2 * * * /home/youruser/apps/tukem/backup.sh >> /home/youruser/apps/tukem/backup.log 2>&1
```

---

## 🐛 Troubleshooting

### Container tidak start

```bash
# Check logs
docker compose -f docker-compose.prod.yml logs

# Check container status
docker ps -a

# Restart specific service
docker compose -f docker-compose.prod.yml restart api
```

### Database connection error

```bash
# Check database logs
docker compose -f docker-compose.prod.yml logs db

# Check database is running
docker compose -f docker-compose.prod.yml exec db pg_isready -U tukem_user

# Test connection
docker compose -f docker-compose.prod.yml exec db psql -U tukem_user -d tukem_db
```

### Port already in use

```bash
# Check what's using the port
sudo netstat -tulpn | grep :8085
sudo netstat -tulpn | grep :3000

# Kill process (ganti PID dengan process ID)
sudo kill -9 PID
```

### Nginx 502 Bad Gateway

```bash
# Check if services are running
docker compose -f docker-compose.prod.yml ps

# Check Nginx error logs
sudo tail -f /var/log/nginx/error.log

# Test backend directly
curl http://localhost:8085/health
```

---

## 📊 Monitoring (Optional)

### Setup Monitoring dengan Docker Stats

```bash
# Real-time stats
docker stats

# Specific container
docker stats tukem-api-prod
```

### Setup Log Rotation

Edit `/etc/logrotate.d/docker-containers`:

```
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size=1M
    missingok
    delaycompress
    copytruncate
}
```

---

## 🔒 Security Best Practices

1. **Change Default Passwords**: Pastikan semua password diubah
2. **Use Strong JWT Secret**: Minimal 32 karakter random
3. **Enable Firewall**: Hanya buka port yang diperlukan
4. **Regular Updates**: Update system dan Docker images secara berkala
5. **SSL Certificate**: Selalu gunakan HTTPS
6. **Backup Regularly**: Setup automated backup
7. **Monitor Logs**: Check logs secara berkala untuk suspicious activity

---

## 📝 Checklist Deployment

- [ ] VPS sudah terinstall Docker & Docker Compose
- [ ] Repository sudah di-clone
- [ ] File `.env` sudah dibuat dengan konfigurasi production
- [ ] Password database sudah diubah
- [ ] JWT secret sudah di-generate
- [ ] Docker Compose production file sudah dibuat
- [ ] Services sudah running (`docker compose ps`)
- [ ] Nginx sudah dikonfigurasi (jika pakai domain)
- [ ] SSL certificate sudah di-setup (jika pakai domain)
- [ ] Firewall sudah dikonfigurasi
- [ ] Backup script sudah dibuat
- [ ] Aplikasi sudah di-test di browser
- [ ] Semua fitur sudah di-verifikasi

---

## 🆘 Support

Jika ada masalah:
1. Check logs: `docker compose -f docker-compose.prod.yml logs`
2. Check container status: `docker compose -f docker-compose.prod.yml ps`
3. Check Nginx logs: `sudo tail -f /var/log/nginx/error.log`
4. Review dokumentasi troubleshooting di atas

---

**Updated:** 21 November 2025

