DATABASE_URL ?= postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable
TEST_DATABASE_URL ?= $(DATABASE_URL)

.PHONY: dev-db stop-db migrate-up migrate-down backend-run backend-test frontend-install frontend-dev frontend-test frontend-typecheck lint test build validate

dev-db:
	docker compose up -d postgres

stop-db:
	docker compose stop postgres

migrate-up:
	docker compose exec -T postgres psql -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.up.sql

migrate-down:
	docker compose exec -T postgres psql -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.down.sql

backend-run:
	set -a; [ -f .env ] && . ./.env; set +a; cd backend && go run ./cmd/api

backend-test:
	cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-mod TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm test

frontend-typecheck:
	cd frontend && npm run typecheck

lint:
	cd backend && gofmt -w .
	cd frontend && npm run typecheck

test: backend-test frontend-test

build:
	cd frontend && npm run build
	cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-mod go build ./cmd/api

validate: lint test build
