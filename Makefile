DATABASE_URL ?= postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable
TEST_DATABASE_URL ?= $(DATABASE_URL)

.PHONY: dev-db stop-db wait-db migrate-up migrate-down reset-dev-db format backend-run backend-test frontend-install frontend-dev frontend-test frontend-typecheck lint test build validate full-validate

dev-db:
	docker compose up -d postgres

stop-db:
	docker compose stop postgres

wait-db:
	for attempt in 1 2 3 4 5 6 7 8 9 10; do docker compose exec -T postgres pg_isready -U marketplace -d marketplace_dev && exit 0; sleep 1; done; docker compose exec -T postgres pg_isready -U marketplace -d marketplace_dev

migrate-up:
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.up.sql
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000002_publish_products.up.sql

migrate-down:
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000002_publish_products.down.sql
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.down.sql

reset-dev-db:
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.down.sql
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000001_init.up.sql
	docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000002_publish_products.up.sql

format:
	cd backend && gofmt -w .

backend-run:
	set -a; [ -f .env ] && . ./.env; set +a; cd backend && go run ./cmd/api

backend-test:
	cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm test

frontend-typecheck:
	cd frontend && npm run typecheck

lint:
	unformatted="$$(gofmt -l backend)"; test -z "$$unformatted" || (printf '%s\n' "$$unformatted"; exit 1)
	cd frontend && npm run typecheck

test: backend-test frontend-test

build:
	cd frontend && npm run build
	cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod go build ./cmd/api

validate: lint test build

full-validate: frontend-install dev-db wait-db reset-dev-db validate
