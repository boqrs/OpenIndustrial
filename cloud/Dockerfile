# ============================================================
# OpenIndustrial Cloud
# Go backend production image
# ============================================================

# ------------------------------------------------------------
# Stage 1: Build
# ------------------------------------------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Copy dependency files first for Docker layer cache
COPY go.mod go.sum ./

RUN go mod download

# Copy backend source
COPY . .

# Build production binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/openindustrial-cloud \
    ./cmd/server


# ------------------------------------------------------------
# Stage 2: Runtime
# ------------------------------------------------------------
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    tzdata

ENV TZ=UTC

COPY --from=builder /app/bin/openindustrial-cloud /app/openindustrial-cloud

# Run as non-root user
RUN addgroup -S openindustrial \
    && adduser -S -G openindustrial openindustrial \
    && chown -R openindustrial:openindustrial /app

USER openindustrial

EXPOSE 8080

ENTRYPOINT ["/app/openindustrial-cloud", "gin"]

