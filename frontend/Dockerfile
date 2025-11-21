FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build application
# This creates .output directory with server/ and public/ subdirectories
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Copy package files to install production dependencies
# Nitro server needs runtime dependencies to serve static files correctly
COPY --from=builder /app/package*.json ./

# Install production dependencies
# Some dependencies are needed at runtime for Nitro to work properly
RUN npm ci --omit=dev

# Copy the entire .output directory from builder
# According to Nuxt docs, we copy .output to root and run from there
# Structure: /app/server/index.mjs and /app/public/_nuxt/...
COPY --from=builder /app/.output ./

# Create symlink for Nitro to find public directory
# Nitro looks for public files in server/chunks/public/ relative to server/index.mjs
# But files are actually in /app/public/, so we create a symlink
RUN mkdir -p /app/server/chunks && ln -s /app/public /app/server/chunks/public

# Set environment variables
# Nitro respects NITRO_PORT/PORT and NITRO_HOST/HOST
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
ENV NITRO_PRESET=node-server

EXPOSE 3000

# Run the Nitro server from the directory containing .output
# Working directory is /app, and .output contents were copied to /app
# So structure is: /app/server/index.mjs and /app/public/_nuxt/...
WORKDIR /app
CMD ["node", "server/index.mjs"]
