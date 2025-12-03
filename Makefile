db-test:
	docker run --name=db-test -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:17.4-alpine -rm

db-it:
	docker exec -it db-test /bin/bash

psql:
	docker exec -it db-test psql -U postgres

migrate:
	goose postgres "user=postgres password=postgres dbname=postgres sslmode=disable" up

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin cmd/bot/main.go

docker-build-n-push:
	docker buildx build --no-cache --platform linux/amd64 -t hack3r11/go-tg-bot:v0.0.1 .
	docker login -u hack3r11 --password-stdin $DOCKER_HUB_TOKEN
	docker push hack3r11/go-tg-bot:v0.0.1