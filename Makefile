SHELL := /bin/sh

COMPOSE ?= docker compose

.PHONY: help up down logs ps build run test fmt tidy openapi migrate-up migrate-down

help:
	@printf '%s\n' \
		'up           Start PostgreSQL, run migrations, and launch the API' \
		'down         Stop containers and remove the database volume' \
		'logs         Follow compose logs' \
		'ps           Show compose services' \
		'build        Build the Go API locally' \
		'run          Run the Go API locally' \
		'test         Run Go tests' \
		'fmt          Format Go sources with gofmt' \
		'tidy         Update Go module metadata' \
		'openapi      Generate openapi/openapi.yaml through TypeSpec' \
		'migrate-up   Apply PostgreSQL migrations through Compose' \
		'migrate-down Roll back PostgreSQL migrations through Compose'

up:
	$(COMPOSE) up -d db migrate-up app

down:
	$(COMPOSE) down -v --remove-orphans

logs:
	$(COMPOSE) logs -f app db

ps:
	$(COMPOSE) ps

build:
	go build ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $$(find cmd internal -name '*.go' -type f)

tidy:
	go mod tidy

openapi:
	$(COMPOSE) run --rm openapi

migrate-up:
	$(COMPOSE) run --rm migrate-up

migrate-down:
	$(COMPOSE) run --rm migrate-down
