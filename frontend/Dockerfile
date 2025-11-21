FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies (needed for nuxt start)
RUN npm ci --omit=dev

# Copy entire .output directory (includes both server and public)
COPY --from=builder /app/.output ./.output
# Copy nuxt config and other necessary files
COPY --from=builder /app/nuxt.config.ts ./nuxt.config.ts

EXPOSE 3000

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

# Use node directly to run the server
CMD ["node", ".output/server/index.mjs"]
