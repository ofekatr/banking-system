include .env

POSTGRES_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable


postgres:
	docker run --name $(POSTGRES_CONTAINER) -p $(POSTGRES_PORT):5432 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --user=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

migrate:
	migrate -path db/migration -database $(POSTGRES_URL) --verbose up

sqlc:
	sqlc generate

dropdb:
	docker exec -it postgres14 dropdb $(POSTGRES_DB)

.PHONY: postgres createdb migrate sqlc dropdb