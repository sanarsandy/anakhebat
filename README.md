# Tukem - Aplikasi Monitoring Tumbuh Kembang Anak

Aplikasi web untuk monitoring tumbuh kembang anak dengan fitur tracking pertumbuhan, milestone perkembangan, dan jadwal imunisasi.

## 🚀 Tech Stack

### Backend
- **Go** (Golang) dengan Echo framework
- **PostgreSQL** database
- **Docker** untuk containerization

### Frontend
- **Vue.js 3** dengan Nuxt 3
- **Tailwind CSS** untuk styling
- **Chart.js** untuk visualisasi grafik
- **Pinia** untuk state management

## 📋 Fitur Utama

- ✅ **Child Profile Management** - Kelola profil anak (multi-profile support)
- ✅ **Growth Tracking** - Tracking berat badan, tinggi badan, lingkar kepala dengan Z-score calculation
- ✅ **Growth Charts** - Visualisasi grafik pertumbuhan dengan WHO standards
- ✅ **Milestone Tracking** - Tracking milestone perkembangan anak
- ✅ **Denver II Assessment** - Assessment perkembangan menggunakan Denver II
- ✅ **Corrected Age Logic** - Perhitungan umur koreksi untuk bayi prematur
- ✅ **PDF Export** - Export laporan lengkap ke PDF
- ✅ **Immunization Schedule** - Jadwal imunisasi berdasarkan IDAI
- ✅ **Recommendations** - Rekomendasi stimulasi berdasarkan milestone

## 🛠️ Setup & Installation

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (untuk development)
- Node.js 18+ (untuk development)

### Quick Start dengan Docker

1. Clone repository:
```bash
git clone https://github.com/sanarsandy/anakhebat.git
cd anakhebat
```

2. Setup environment variables:
```bash
# Copy .env.example ke .env (jika ada)
# Atau set environment variables di docker-compose.yml
```

3. Start services:
```bash
docker-compose up -d
```

4. Akses aplikasi:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

### Development Setup

#### Backend
```bash
cd backend
go mod download
go run main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

## 📁 Project Structure

```
tukem/
├── backend/              # Go backend
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── middleware/      # Middleware (JWT, etc)
│   ├── migrations/      # Database migrations
│   ├── utils/           # Utilities (Z-score, age calc, etc)
│   └── main.go          # Entry point
├── frontend/            # Nuxt 3 frontend
│   ├── components/      # Vue components
│   ├── pages/           # Pages/routes
│   ├── stores/          # Pinia stores
│   ├── composables/     # Composables
│   └── nuxt.config.ts   # Nuxt config
├── docker-compose.yml   # Docker setup
└── README.md           # This file
```

## 🔐 Environment Variables

Lihat `.env.example` untuk daftar environment variables yang diperlukan.

**Important:** Jangan commit file `.env` ke repository. Gunakan `.env.example` sebagai template.

## 📊 Database

Database menggunakan PostgreSQL. Migrations ada di `backend/migrations/`.

Untuk apply migrations, jalankan:
```bash
# Migrations di-apply otomatis saat backend start
# Atau manual via Go code
```

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm run test
```

## 📝 API Documentation

API endpoints tersedia di:
- Base URL: `http://localhost:8080/api`
- Health check: `GET /health`

### Main Endpoints:
- `POST /api/auth/register` - Register user
- `POST /api/auth/login` - Login
- `GET /api/children` - Get children list
- `POST /api/children` - Create child
- `GET /api/children/:id/measurements` - Get measurements
- `POST /api/children/:id/measurements` - Create measurement
- `GET /api/children/:id/export-pdf` - Export PDF report
- Dan lainnya...

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is private/proprietary.

## 👥 Authors

- **Sanar Sandy** - [@sanarsandy](https://github.com/sanarsandy)

## 🙏 Acknowledgments

- WHO Growth Standards
- IDAI Immunization Schedule
- Denver II Developmental Screening

---

**Note:** Project ini masih dalam tahap pengembangan aktif.

