DATABASE_URL ?= postgres://medicabinet:medicabinet@localhost:5432/medicabinet?sslmode=disable

.PHONY: up down logs sqlc migrate-up test dev-server dev-client

## up: start Postgres + backend via docker compose
up:
	docker compose up --build

## down: stop and remove compose services
down:
	docker compose down

## logs: tail compose logs
logs:
	docker compose logs -f

## sqlc: regenerate type-safe DB code (requires queries in server/db/queries)
sqlc:
	cd server && sqlc generate

## migrate-up: apply migrations (requires golang-migrate CLI)
migrate-up:
	migrate -path server/db/migrations -database "$(DATABASE_URL)" up

## test: run the Go test suite
test:
	cd server && go test ./...

## dev-server: run the Go server locally
dev-server:
	cd server && go run ./cmd/server

## dev-client: run the Next.js dev server locally
dev-client:
	cd client && npm run dev
