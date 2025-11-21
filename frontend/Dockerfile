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

# Create symlink to ensure server can find public directory
# Server looks for public in .output/public, but we copied .output to /app
# So public is now at /app/public, which is correct
RUN ls -la /app/ && ls -la /app/public/ 2>/dev/null || echo "Public directory check"

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
