# Troubleshooting Deployment Issues

## Error: Port Already in Use

### Problem
```
Error response from daemon: failed to bind host port for 127.0.0.1:8085: address already in use
```

### Solution 1: Check What's Using the Port

```bash
# Check what's using port 8085
sudo netstat -tulpn | grep :8085
# atau
sudo lsof -i :8085
# atau
sudo ss -tulpn | grep :8085
```

### Solution 2: Kill the Process

```bash
# Jika ada process yang menggunakan port 8085, kill process tersebut
# Ganti PID dengan process ID dari command di atas
sudo kill -9 PID

# Atau kill semua process di port 8085
sudo fuser -k 8085/tcp
```

### Solution 3: Stop Existing Containers

```bash
# Check jika ada container Docker yang masih running
docker ps -a | grep 8080

# Stop semua container Tukem yang mungkin masih running
docker stop tukem-api-prod tukem-api 2>/dev/null
docker rm tukem-api-prod tukem-api 2>/dev/null

# Atau stop semua container
docker compose -f docker-compose.prod.yml down
docker compose down 2>/dev/null
```

### Solution 4: Use Different Port

Edit `docker-compose.prod.yml` dan ubah port:

```yaml
api:
  ports:
    - "127.0.0.1:8086:8080"  # Ubah dari 8085 ke 8086
```

Dan update Nginx config untuk point ke port 8086:

```nginx
location /api {
    proxy_pass http://127.0.0.1:8086;  # Ubah ke 8086
    ...
}
```

---

## Error: Port 3000 Already in Use

### Solution

```bash
# Check port 3000
sudo netstat -tulpn | grep :3000

# Kill process
sudo fuser -k 3000/tcp

# Atau stop existing containers
docker stop tukem-app-prod tukem-app 2>/dev/null
docker rm tukem-app-prod tukem-app 2>/dev/null
```

---

## Error: Database Connection Failed

### Check Database Container

```bash
# Check if database is running
docker compose -f docker-compose.prod.yml ps db

# Check database logs
docker compose -f docker-compose.prod.yml logs db

# Test database connection
docker compose -f docker-compose.prod.yml exec db pg_isready -U tukem_user
```

### Solution

```bash
# Restart database
docker compose -f docker-compose.prod.yml restart db

# Wait for database to be healthy
docker compose -f docker-compose.prod.yml ps db
```

---

## Error: Container Keeps Restarting

### Check Logs

```bash
# Check logs for errors
docker compose -f docker-compose.prod.yml logs api
docker compose -f docker-compose.prod.yml logs app

# Check container status
docker compose -f docker-compose.prod.yml ps
```

### Common Causes

1. **Environment variables missing**: Check `.env` file
2. **Database not ready**: Wait for database to be healthy
3. **Port conflict**: Check port usage
4. **Build errors**: Rebuild containers

### Solution

```bash
# Rebuild containers
docker compose -f docker-compose.prod.yml build --no-cache

# Start fresh
docker compose -f docker-compose.prod.yml down -v
docker compose -f docker-compose.prod.yml up -d
```

---

## Error: Permission Denied

### Solution

```bash
# Fix permissions for backup script
chmod +x backup.sh

# Fix permissions for Docker socket (if needed)
sudo chmod 666 /var/run/docker.sock
# (Not recommended for production, better to add user to docker group)
sudo usermod -aG docker $USER
# Logout and login again
```

---

## Clean Up Everything and Start Fresh

```bash
# Stop and remove all containers
docker compose -f docker-compose.prod.yml down -v
docker compose down -v 2>/dev/null

# Remove all Tukem containers
docker ps -a | grep tukem | awk '{print $1}' | xargs docker rm -f

# Remove volumes (⚠️ Will delete data!)
docker volume ls | grep tukem | awk '{print $2}' | xargs docker volume rm

# Clean up networks
docker network prune -f

# Start fresh
docker compose -f docker-compose.prod.yml up -d
```

---

## Quick Fix Script

Buat file `fix-deployment.sh`:

```bash
#!/bin/bash

echo "Stopping all Tukem containers..."
docker compose -f docker-compose.prod.yml down 2>/dev/null
docker compose down 2>/dev/null

echo "Killing processes on ports 8085 and 3000..."
sudo fuser -k 8085/tcp 2>/dev/null
sudo fuser -k 3000/tcp 2>/dev/null

echo "Cleaning up..."
docker network prune -f

echo "Starting services..."
docker compose -f docker-compose.prod.yml up -d

echo "Checking status..."
docker compose -f docker-compose.prod.yml ps
```

Run:
```bash
chmod +x fix-deployment.sh
./fix-deployment.sh
```

