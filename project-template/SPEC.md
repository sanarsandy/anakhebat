# Spesifikasi Teknis Project Template
## Full-Stack Web Application dengan Go + Nuxt.js

---

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Arsitektur Sistem](#arsitektur-sistem)
3. [Tech Stack](#tech-stack)
4. [Struktur Project](#struktur-project)
5. [Database Schema](#database-schema)
6. [API Endpoints](#api-endpoints)
7. [Authentication & Authorization](#authentication--authorization)
8. [Deployment](#deployment)
9. [Environment Variables](#environment-variables)
10. [Development Workflow](#development-workflow)

---

## 🎯 Overview

Template ini adalah full-stack web application dengan arsitektur:
- **Backend**: Go (Echo Framework) + PostgreSQL
- **Frontend**: Nuxt.js 3 (Vue.js 3) + TailwindCSS
- **Containerization**: Docker & Docker Compose
- **Authentication**: JWT + Google OAuth 2.0
- **State Management**: Pinia

### Fitur Utama

- ✅ RESTful API dengan Go Echo Framework
- ✅ JWT Authentication & Authorization
- ✅ Google OAuth 2.0 Integration
- ✅ PostgreSQL Database dengan Migration System
- ✅ Nuxt.js 3 dengan SSR/SSG Support
- ✅ Docker & Docker Compose untuk Development & Production
- ✅ Multi-stage Docker Build untuk Optimasi
- ✅ Environment-based Configuration
- ✅ CORS Configuration
- ✅ Database Seeding System

---

## 🏗️ Arsitektur Sistem

```
┌─────────────────┐
│   Nginx (Prod)  │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼────┐
│ Nuxt  │ │  Go   │
│  App  │ │  API  │
└───┬───┘ └───┬───┘
    │         │
    └────┬────┘
         │
    ┌────▼────┐
    │PostgreSQL│
    └─────────┘
```

### Komponen Utama

1. **Frontend (Nuxt.js)**
   - Port: 3000 (Development), 3000 (Production)
   - SSR dengan Nitro Server
   - Client-side routing dengan Vue Router
   - State management dengan Pinia

2. **Backend (Go Echo)**
   - Port: 8080 (Development), 8085 (Production)
   - RESTful API
   - JWT Middleware untuk protected routes
   - Database connection pooling

3. **Database (PostgreSQL)**
   - Port: 5432
   - Version: 15-alpine
   - Persistent volume untuk data

---

## 🛠️ Tech Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.20+ | Programming Language |
| Echo Framework | v4.11+ | Web Framework |
| PostgreSQL | 15-alpine | Database |
| JWT | v5.2.0 | Authentication |
| OAuth2 | v0.8.0 | Google OAuth |
| Bcrypt | Latest | Password Hashing |
| lib/pq | v1.10+ | PostgreSQL Driver |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| Nuxt.js | 3.10+ | Framework |
| Vue.js | 3.4+ | UI Framework |
| Pinia | 2.1+ | State Management |
| TailwindCSS | 3.3+ | CSS Framework |
| Chart.js | 4.5+ | Data Visualization |
| Vue Chart.js | 5.3+ | Chart Components |

### DevOps

| Technology | Purpose |
|------------|---------|
| Docker | Containerization |
| Docker Compose | Multi-container Orchestration |
| Nginx | Reverse Proxy (Production) |

---

## 📁 Struktur Project

```
project-root/
├── backend/                    # Go Backend
│   ├── cmd/                    # CLI commands
│   ├── data/                   # Seed data (JSON files)
│   ├── db/                     # Database connection
│   │   └── db.go
│   ├── handlers/               # HTTP handlers
│   │   ├── auth.go
│   │   ├── google_auth.go
│   │   └── ...
│   ├── middleware/             # HTTP middleware
│   │   └── jwt.go
│   ├── migrations/             # Database migrations
│   │   ├── 001_init_schema.sql
│   │   └── ...
│   ├── models/                 # Data models
│   │   ├── user.go
│   │   └── ...
│   ├── utils/                  # Utility functions
│   │   ├── seed_*.go
│   │   └── ...
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── init_db.sql            # Initial database setup
│   └── main.go
│
├── frontend/                   # Nuxt.js Frontend
│   ├── assets/                 # Static assets
│   │   └── css/
│   ├── components/             # Vue components
│   ├── composables/            # Composable functions
│   ├── layouts/                # Layout components
│   ├── middleware/             # Route middleware
│   ├── pages/                  # Pages (auto-routing)
│   ├── plugins/                # Nuxt plugins
│   ├── stores/                 # Pinia stores
│   ├── Dockerfile
│   ├── nuxt.config.ts
│   ├── package.json
│   └── tailwind.config.js
│
├── docs/                       # Documentation (gitignored)
├── project-template/           # Project template
│   ├── SPEC.md                 # This file
│   ├── README.md               # Template usage guide
│   ├── template/               # Template files
│   └── scripts/                # Generation scripts
│
├── docker-compose.yml          # Development setup
├── docker-compose.prod.yml     # Production setup
├── .env.example                # Environment variables template
├── .gitignore
└── README.md
```

---

## 🗄️ Database Schema

### Core Tables

#### `users`
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),  -- NULL for Google-only users
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'parent',
    google_id VARCHAR(255) UNIQUE,
    auth_provider VARCHAR(50) DEFAULT 'email',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX idx_users_email_auth_provider 
    ON users (email, auth_provider) 
    WHERE auth_provider IN ('email', 'both');
CREATE UNIQUE INDEX idx_users_email_google 
    ON users (email) 
    WHERE auth_provider = 'google';
```

#### Migration System

Migrations disimpan di `backend/migrations/` dengan format:
- `001_*.sql` - Initial schema
- `002_*.sql` - Subsequent changes
- Format: `{number}_{description}.sql`

**Best Practices:**
- Setiap migration harus idempotent (gunakan `IF NOT EXISTS`)
- Jangan hapus migration yang sudah di-deploy
- Test migration di development sebelum production

---

## 🔌 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | No |
| POST | `/api/auth/login` | Login with email/password | No |
| GET | `/api/auth/google` | Get Google OAuth URL | No |
| GET | `/api/auth/google/callback` | Google OAuth callback | No |

### Protected Routes (Require JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/me` | Get current user info |
| GET | `/api/children` | Get user's children |
| POST | `/api/children` | Create new child |
| GET | `/api/children/:id` | Get child details |
| PUT | `/api/children/:id` | Update child |
| DELETE | `/api/children/:id` | Delete child |

### JWT Token Format

```json
{
  "user_id": "uuid",
  "role": "parent|admin",
  "exp": 1234567890
}
```

**Token Lifetime**: 72 hours (3 days)

---

## 🔐 Authentication & Authorization

### JWT Authentication Flow

1. **Email/Password Login**
   ```
   POST /api/auth/login
   → Verify credentials
   → Generate JWT token
   → Return token + user data
   ```

2. **Google OAuth Flow**
   ```
   GET /api/auth/google
   → Redirect to Google
   → User authorizes
   → GET /api/auth/google/callback?code=...
   → Exchange code for token
   → Find/create user
   → Generate JWT token
   → Redirect to frontend with token
   ```

3. **Protected Routes**
   ```
   Request with Authorization: Bearer <token>
   → JWT Middleware validates token
   → Extract user_id from claims
   → Attach to context
   → Handler processes request
   ```

### Auth Provider Types

- `email`: Email/password only
- `google`: Google OAuth only
- `both`: Both methods (future feature)

**Important**: Email dapat digunakan untuk multiple auth providers dengan unique indexes yang berbeda.

---

## 🚀 Deployment

### Development

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

### Production

```bash
# Build and start
docker compose -f docker-compose.prod.yml up -d --build

# View logs
docker compose -f docker-compose.prod.yml logs -f api

# Restart service
docker compose -f docker-compose.prod.yml restart api
```

### Nginx Configuration (Production)

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Proxy all requests to Nuxt app
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

---

## 🔧 Environment Variables

### Backend (.env)

```env
# Database
DB_HOST=db
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432

# JWT
JWT_SECRET=your-secure-jwt-secret-here

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback

# Frontend
FRONTEND_URL=https://yourdomain.com

# Server
PORT=8080
ENV=production
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

### Frontend (.env)

```env
NUXT_PUBLIC_API_BASE=https://yourdomain.com/api
NODE_ENV=production
```

**Security Notes:**
- ❌ Jangan commit `.env` file ke git
- ✅ Gunakan `.env.example` sebagai template
- ✅ Set environment variables di production server
- ✅ Gunakan secrets management untuk production

---

## 💻 Development Workflow

### 1. Setup Project

```bash
# Clone template
git clone <template-repo> my-project
cd my-project

# Copy environment file
cp .env.example .env
# Edit .env with your values

# Start services
docker compose up -d

# Run migrations (if needed)
docker compose exec api sh -c "cd /app && go run main.go"
```

### 2. Backend Development

```bash
# Enter backend container
docker compose exec api sh

# Run migrations manually
psql -U $DB_USER -d $DB_NAME -f migrations/001_init_schema.sql

# Test API
curl http://localhost:8080/health
```

### 3. Frontend Development

```bash
# Frontend auto-reloads on file changes
# Access at http://localhost:3000
```

### 4. Database Migrations

```bash
# Create new migration
touch backend/migrations/009_new_feature.sql

# Write migration (idempotent)
# Example:
# ALTER TABLE users ADD COLUMN IF NOT EXISTS new_field VARCHAR(255);

# Apply migration
docker compose exec db psql -U $DB_USER -d $DB_NAME -f /path/to/migration.sql
```

### 5. Adding New Features

1. **New API Endpoint**
   - Add handler in `backend/handlers/`
   - Register route in `backend/main.go`
   - Add model if needed in `backend/models/`

2. **New Frontend Page**
   - Create file in `frontend/pages/`
   - Nuxt auto-generates route
   - Use Pinia store for state management

3. **New Database Table**
   - Create migration in `backend/migrations/`
   - Add model in `backend/models/`
   - Update handlers as needed

---

## 📝 Best Practices

### Code Organization

- **Handlers**: Keep handlers focused on HTTP concerns
- **Models**: Define data structures and validation
- **Utils**: Reusable business logic
- **Middleware**: Cross-cutting concerns (auth, logging, etc.)

### Error Handling

```go
// Good
if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{
        "error": "Descriptive error message",
    })
}

// Bad
if err != nil {
    return err
}
```

### Security

- ✅ Always validate input
- ✅ Use parameterized queries (prevent SQL injection)
- ✅ Hash passwords with bcrypt
- ✅ Use HTTPS in production
- ✅ Validate JWT tokens in middleware
- ✅ Set secure CORS policies

### Performance

- ✅ Use database indexes
- ✅ Implement connection pooling
- ✅ Cache frequently accessed data
- ✅ Optimize Docker images (multi-stage builds)
- ✅ Use CDN for static assets (production)

---

## 🧪 Testing

### Manual Testing

```bash
# Test API health
curl http://localhost:8080/health

# Test authentication
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Test protected route
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer <token>"
```

### Database Testing

```bash
# Connect to database
docker compose exec db psql -U $DB_USER -d $DB_NAME

# Run queries
SELECT * FROM users;
```

---

## 📚 Additional Resources

- [Go Echo Documentation](https://echo.labstack.com/)
- [Nuxt.js Documentation](https://nuxt.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [JWT.io](https://jwt.io/) - JWT Debugger

---

## 🔄 Version History

- **v1.0.0** (2025-01-XX)
  - Initial template
  - Go + Nuxt.js stack
  - JWT + Google OAuth
  - Docker setup
  - Migration system

---

## 📄 License

This template is provided as-is for use in your projects.

---

**Last Updated**: 2025-01-XX
**Maintained By**: Your Team

