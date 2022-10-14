-include app.env

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable


postgres:
	docker run --name postgres14 -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --user=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

migrateup:
	migrate -path db/migrations -database $(DB_URL) --verbose up

migratedown:
	migrate -path db/migrations -database $(DB_URL) --verbose down

sqlc:
	sqlc generate

dropdb:
	docker exec -it postgres14 dropdb $(DB_NAME)

test:
	go test -v -cover ./...

server:
	go run main.go

# .PHONY: postgres createdb migrate sqlc dropdb