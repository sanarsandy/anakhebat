# Fix Database Container Not Running

## Check Status Dulu

```bash
# Check container status
docker ps -a | grep tukem-db-prod

# Check logs
docker logs tukem-db-prod

# Check health
docker inspect tukem-db-prod --format='{{.State.Health.Status}}'
```

## Solusi 1: Fresh Start dengan Default Values

```bash
# Stop semua
docker compose -f docker-compose.prod.yml down

# Hapus volume
docker volume rm tukem_postgres_data 2>/dev/null || true

# Buat .env file dengan default values
cat > .env << EOF
DB_HOST=db
DB_USER=tukem_user
DB_PASSWORD=tukem_password
DB_NAME=tukem_db
DB_PORT=5432
PORT=8080
ENV=production
JWT_SECRET=$(openssl rand -base64 32)
FRONTEND_URL=http://localhost:3000
NUXT_PUBLIC_API_BASE=http://localhost:8080/api
CORS_ALLOWED_ORIGINS=http://localhost:3000
EOF

# Start
docker compose -f docker-compose.prod.yml up -d db

# Check logs
docker compose -f docker-compose.prod.yml logs -f db
```

## Solusi 2: Run Database Manual (Tanpa Docker Compose)

```bash
# Stop container
docker stop tukem-db-prod 2>/dev/null || true
docker rm tukem-db-prod 2>/dev/null || true

# Run manual
docker run -d \
  --name tukem-db-prod \
  --restart unless-stopped \
  -e POSTGRES_USER=tukem_user \
  -e POSTGRES_PASSWORD=tukem_password \
  -e POSTGRES_DB=tukem_db \
  -v tukem_postgres_data:/var/lib/postgresql/data \
  -p 127.0.0.1:5432:5432 \
  postgres:15-alpine

# Check logs
docker logs -f tukem-db-prod
```

## Solusi 3: Simplify Docker Compose

Edit `docker-compose.prod.yml`, ubah bagian db:

```yaml
db:
  image: postgres:15-alpine
  container_name: tukem-db-prod
  restart: unless-stopped
  environment:
    POSTGRES_USER: tukem_user
    POSTGRES_PASSWORD: tukem_password
    POSTGRES_DB: tukem_db
  volumes:
    - postgres_data:/var/lib/postgresql/data
  networks:
    - tukem-network
  # Remove healthcheck untuk sementara
```

Lalu:
```bash
docker compose -f docker-compose.prod.yml up -d db
```

## Check Error Spesifik

```bash
# Run check script
chmod +x check-db.sh
./check-db.sh
```

