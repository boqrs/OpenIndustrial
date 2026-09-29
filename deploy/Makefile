# ============================================================
# OpenIndustrial Deploy Makefile
# ============================================================

SHELL := /bin/bash

COMPOSE := docker compose
COMPOSE_FILE := docker-compose.yml

# ============================================================
# Default
# ============================================================

.PHONY: help
help:
	@echo ""
	@echo "OpenIndustrial Deploy"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build Cloud backend image"
	@echo "  make run         Start all Docker Compose services"
	@echo "  make stop        Stop all Docker Compose services"
	@echo "  make restart     Restart all Docker Compose services"
	@echo "  make ps          Show Docker Compose services"
	@echo "  make logs        Show Docker Compose logs"
	@echo "  make cloud-logs  Show Cloud backend logs"
	@echo "  make nginx-logs  Show Nginx logs"
	@echo "  make db-logs     Show PostgreSQL logs"
	@echo "  make config      Validate Docker Compose configuration"
	@echo "  make clean       Stop services and remove containers"
	@echo ""


# ============================================================
# Build
# ============================================================

.PHONY: build
build:
	@echo "==> Building OpenIndustrial Cloud..."
	docker build \
		-t openindustrial-cloud:latest \
		../cloud


# ============================================================
# Run
# ============================================================

.PHONY: run
run:
	@echo "==> Starting OpenIndustrial..."
	$(COMPOSE) -f $(COMPOSE_FILE) up -d


# ============================================================
# Stop
# ============================================================

.PHONY: stop
stop:
	@echo "==> Stopping OpenIndustrial..."
	$(COMPOSE) -f $(COMPOSE_FILE) stop


# ============================================================
# Restart
# ============================================================

.PHONY: restart
restart:
	@echo "==> Restarting OpenIndustrial..."
	$(COMPOSE) -f $(COMPOSE_FILE) restart


# ============================================================
# Status
# ============================================================

.PHONY: ps
ps:
	$(COMPOSE) -f $(COMPOSE_FILE) ps


# ============================================================
# Logs
# ============================================================

.PHONY: logs
logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f


.PHONY: cloud-logs
cloud-logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f cloud


.PHONY: nginx-logs
nginx-logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f nginx


.PHONY: db-logs
db-logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f postgres


# ============================================================
# Validate
# ============================================================

.PHONY: config
config:
	$(COMPOSE) -f $(COMPOSE_FILE) config


# ============================================================
# Clean
# ============================================================

.PHONY: clean
clean:
	@echo "==> Stopping and removing OpenIndustrial containers..."
	$(COMPOSE) -f $(COMPOSE_FILE) down
