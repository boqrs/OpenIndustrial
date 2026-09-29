# ============================================================
# OpenIndustrial Console
# Vue frontend production image
# ============================================================

# ------------------------------------------------------------
# Stage 1: Build frontend
# ------------------------------------------------------------
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files first for Docker layer cache
COPY package.json package-lock.json ./

# Install exactly according to package-lock.json
RUN npm ci

# Copy frontend source
COPY . .

# Build production frontend
RUN npm run build


# ------------------------------------------------------------
# Stage 2: Nginx runtime
# ------------------------------------------------------------
FROM nginx:1.29-alpine

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy Vue production build
COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx configuration will be supplied by docker-compose.
#
# The reason we don't COPY nginx.conf here is that the
# reverse-proxy topology belongs to the deployment layer.
#
# Expected deployment configuration:
#
#   /api/* -> cloud:8080
#   /*     -> Vue SPA
#
# See:
#   deploy/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
