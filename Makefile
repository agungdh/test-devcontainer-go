.PHONY: run migrate-up migrate-down tidy
include .env
export

run:
	go run ./cmd/api

migrate-up:
	goose -dir db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	goose -dir db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

tidy:
	go mod tidy && go vet ./...
