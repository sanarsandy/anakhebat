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

# Copy only the .output directory from builder
# According to Nuxt/Nitro docs, .output is fully standalone
# and contains everything needed to run the server
COPY --from=builder /app/.output ./

# Set environment variables
# Nitro respects NITRO_PORT/PORT and NITRO_HOST/HOST
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

EXPOSE 3000

# Run the Nitro server
# According to official docs: node .output/server/index.mjs
# Since we copied .output to /app, the path is server/index.mjs
CMD ["node", "server/index.mjs"]
