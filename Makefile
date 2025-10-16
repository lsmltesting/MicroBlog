include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@db:$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR=./db/migrations


migrate-up:
	docker-compose run --rm migrator /go/bin/goose -dir /migrations postgres "$(DB_URL)" up

migrate-down:
	docker-compose run --rm migrator /go/bin/goose -dir /migrations postgres "$(DB_URL)" down

migrate-new:
	docker-compose run --rm migrator /go/bin/goose -dir /migrations create $(name) sql

rebuild:
	docker-compose down -v
	docker-compose up -d db
	sleep 10
	make migrate-up
