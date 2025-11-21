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

# Copy .output directory from builder
# Keep the .output structure as server expects it
COPY --from=builder /app/.output ./.output

# Set environment variables
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000
ENV NITRO_PRESET=node-server

EXPOSE 3000

# Run production server from .output directory
# Server expects to run from directory containing .output
WORKDIR /app
CMD ["node", ".output/server/index.mjs"]
