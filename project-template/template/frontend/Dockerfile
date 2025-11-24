FROM node:18-alpine AS base

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

FROM base AS development

FROM base AS builder
RUN npm run build

FROM node:18-alpine
WORKDIR /app

COPY --from=builder /app/.output ./
RUN mkdir -p /app/server/chunks && ln -s /app/public /app/server/chunks/public

ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
ENV NITRO_PRESET=node-server

EXPOSE 3000

WORKDIR /app
CMD ["node", "server/index.mjs"]

