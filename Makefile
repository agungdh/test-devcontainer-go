.PHONY: run migrate-up migrate-down db-reset tidy

# NOTE: jangan `include .env` di sini — itu nimpa env asli container
# (DB_HOST=postgres dari compose) dengan isi file. `go run` sudah
# auto-load .env via godotenv (tidak nimpa env yang sudah ada).
# Target goose pakai fallback ${VAR:-default} supaya env asli menang.
DSN = postgres://$${DB_USER:-postgres}:$${DB_PASSWORD:-postgres}@$${DB_HOST:-postgres}:$${DB_PORT:-5432}/$${DB_NAME:-be_app_dev}?sslmode=disable

run:
	go run ./cmd/api

migrate-up:
	goose -dir db/migrations postgres "$(DSN)" up

migrate-down:
	goose -dir db/migrations postgres "$(DSN)" down

# Reset dev DB: rollback SEMUA migrasi lalu apply ulang dari nol.
# DESTRUKTIF — semua data hilang. Jangan dipakai di production.
db-reset:
	@echo "WARNING: reset database, semua data hilang..."
	goose -dir db/migrations postgres "$(DSN)" reset
	goose -dir db/migrations postgres "$(DSN)" up

tidy:
	go mod tidy && go vet ./...
