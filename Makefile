.PHONY:
.SILENT:

run:
	go run ./cmd/app/main.go

rundb:
	docker run --name='$(DB_CONTAINER_NAME)' \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		--rm -d postgres