-include .env
export

build-image:
	docker build -t app:latest .

postgres:
	docker run --name postgres15-server -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:15-alpine \
	|| docker start postgres15-server

compose-up:
	docker compose up --build -d && docker compose logs -f

compose-down:
	docker compose down --remove-orphans

create-db:
	docker exec -it postgres15-server createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

drop-db:
	docker exec -it postgres15-server dropdb $(DB_NAME)

migrate-up:
	migrate -verbose -path sql/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDR):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -verbose -path sql/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDR):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

sqlc:
	env \
		DB_USER=$(DB_USER)\
		DB_PASSWORD=$(DB_PASSWORD)\
		DB_ADDR=$(DB_ADDR)\
		DB_PORT=$(DB_PORT)\
		DB_NAME=$(DB_NAME) \
		sqlc -f sql/sqlc.yaml vet && sqlc -f sql/sqlc.yaml generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc