.PHONY: run build migrate tidy swagger

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/api ./cmd/api/main.go

migrate:
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/001_create_users.sql
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/002_create_books.sql

tidy:
	go mod tidy

swagger:
	swag init -g cmd/api/main.go -o docs