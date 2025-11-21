FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Copy .output directory from builder to /app
# This includes both server/ and public/ directories
COPY --from=builder /app/.output ./

# Also copy public directory separately for volume mount
# This allows Nginx to serve static files directly
RUN mkdir -p /app/public && cp -r /app/.output/public/* /app/public/ 2>/dev/null || true

# Set environment variables
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000
ENV NITRO_PRESET=node-server

EXPOSE 3000

# Run production server
# Path is /app/server/index.mjs because .output was copied to /app
WORKDIR /app
CMD ["node", "server/index.mjs"]
