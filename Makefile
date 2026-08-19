.PHONY: run rundb stopdb
.SILENT:

include .env
export

run:
	go run ./cmd/app/main.go

rundb:
	docker run --name='$(DB_CONTAINER_NAME)' \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgres

stopdb:
	docker stop $(DB_CONTAINER_NAME)
	docker rm $(DB_CONTAINER_NAME)

migrate:
	go run ./migrations/migrate.go