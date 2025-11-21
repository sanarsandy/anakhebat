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

# Copy entire .output directory from builder
# According to Nuxt docs, copy .output to root and run from there
COPY --from=builder /app/.output ./

# Set environment variables
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

# Run production server
# Nuxt docs: run "node server/index.mjs" from directory containing .output
# Since we copied .output to /app, the structure is:
# /app/server/index.mjs
# /app/public/_nuxt/...
CMD ["node", "server/index.mjs"]
