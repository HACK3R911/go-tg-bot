.PHONY: goose migrate status create

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable


migrate-db:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable goose up -dir ./migrations

status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable goose status -dir ./migrations

create:
	goose create $(name) sql

rollback:
	goose down

db-test:
	sudo docker run --rm --name=db-test -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:17.4-alpine

db-it:
	docker exec -it db-test /bin/bash

psql:
	docker exec -it db-test psql -U postgres

migrate:
	goose postgres "user=postgres password=postgres dbname=postgres sslmode=disable" up

build:
	@echo "Building..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin ./cmd/bot
	@echo "Build complete"

docker-build-n-push:
	docker buildx build --no-cache --platform linux/amd64 -t hack3r11/go-tg-bot:v0.0.1 .
	docker login -u hack3r11 --password-stdin $DOCKER_HUB_TOKEN
	docker push hack3r11/go-tg-bot:v0.0.1