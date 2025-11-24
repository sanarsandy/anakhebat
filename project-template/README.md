# Project Template - Full Stack Go + Nuxt.js

Template project untuk membuat aplikasi full-stack dengan Go (Echo) backend dan Nuxt.js frontend.

## 🚀 Quick Start

### Menggunakan Template

1. **Generate Project Baru**

```bash
# Jalankan script generator
cd project-template/scripts
./generate-project.sh my-new-project

# Atau manual copy
cp -r template ../my-new-project
cd ../my-new-project
```

2. **Setup Environment**

```bash
# Copy environment template
cp .env.example .env

# Edit .env dengan konfigurasi Anda
nano .env
```

3. **Start Development**

```bash
# Start semua services
docker compose up -d

# Check logs
docker compose logs -f

# Access aplikasi
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## 📋 Prerequisites

- Docker & Docker Compose
- Git
- (Optional) Go 1.20+ untuk development lokal
- (Optional) Node.js 18+ untuk development lokal

## 📁 Struktur Template

```
template/
├── backend/              # Go backend
├── frontend/             # Nuxt.js frontend
├── docker-compose.yml    # Development setup
├── docker-compose.prod.yml  # Production setup
├── .env.example          # Environment variables
└── .gitignore
```

## 🔧 Customization

### 1. Ganti Nama Project

```bash
# Ganti semua referensi "tukem" atau nama project lama
# dengan nama project baru Anda

# Backend
sed -i 's/tukem-backend/my-project-backend/g' backend/go.mod
sed -i 's/tukem/my-project/g' backend/main.go

# Frontend
sed -i 's/tukem-frontend/my-project-frontend/g' frontend/package.json

# Docker Compose
sed -i 's/tukem/my-project/g' docker-compose.yml
```

### 2. Setup Database

```bash
# Edit init_db.sql dengan schema Anda
nano backend/init_db.sql

# Atau buat migration baru
touch backend/migrations/001_init_schema.sql
```

### 3. Setup Google OAuth (Optional)

1. Buat project di [Google Cloud Console](https://console.cloud.google.com/)
2. Enable Google+ API
3. Create OAuth 2.0 credentials
4. Set authorized redirect URIs
5. Update `.env` dengan credentials

### 4. Customize Frontend

```bash
# Edit nuxt.config.ts
nano frontend/nuxt.config.ts

# Edit app.vue untuk layout utama
nano frontend/app.vue

# Tambah pages baru
touch frontend/pages/my-page.vue
```

## 🗄️ Database Setup

### Development

```bash
# Database akan otomatis dibuat saat docker compose up
# Untuk run migrations manual:

docker compose exec db psql -U tukem_user -d tukem_db -f /path/to/migration.sql
```

### Production

```bash
# Setup database
docker compose -f docker-compose.prod.yml exec db psql -U $DB_USER -d $DB_NAME -f /path/to/migration.sql

# Backup database
docker compose -f docker-compose.prod.yml exec db pg_dump -U $DB_USER $DB_NAME > backup.sql
```

## 🔐 Security Checklist

- [ ] Ganti `JWT_SECRET` dengan random string yang kuat
- [ ] Ganti database passwords
- [ ] Setup HTTPS di production
- [ ] Konfigurasi CORS dengan benar
- [ ] Review dan update `.gitignore`
- [ ] Jangan commit `.env` file
- [ ] Setup firewall rules
- [ ] Enable database SSL (production)

## 📦 Deployment

### VPS Deployment

1. **Clone Project**

```bash
git clone <your-repo> /var/www/my-project
cd /var/www/my-project
```

2. **Setup Environment**

```bash
cp .env.example .env
nano .env  # Edit dengan production values
```

3. **Build & Start**

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

4. **Setup Nginx**

```bash
# Copy nginx config
sudo cp nginx-config.conf /etc/nginx/sites-available/my-project
sudo ln -s /etc/nginx/sites-available/my-project /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

5. **Setup SSL (Let's Encrypt)**

```bash
sudo certbot --nginx -d yourdomain.com
```

## 🧪 Testing

### Test API

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","full_name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Test Frontend

```bash
# Access di browser
http://localhost:3000
```

## 📚 Documentation

- [Spesifikasi Teknis Lengkap](./SPEC.md)
- [Go Echo Docs](https://echo.labstack.com/)
- [Nuxt.js Docs](https://nuxt.com/)

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :8080
lsof -i :3000

# Kill process or change port in docker-compose.yml
```

### Database Connection Error

```bash
# Check database is running
docker compose ps

# Check database logs
docker compose logs db

# Test connection
docker compose exec db psql -U $DB_USER -d $DB_NAME
```

### Frontend Build Error

```bash
# Clear node_modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
```

## 🤝 Contributing

1. Fork the template
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

## 📄 License

This template is provided as-is. Use it freely in your projects.

---

**Happy Coding! 🚀**

