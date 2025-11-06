build:
	GOOS=linux GOARCH=amd64 go build -o ./bin cmd/main.go

docker-build-n-push:
	docker buildx build --no-cache --platform linux/amd64 -t <REGESTRY>/go-tg-bot:v0.0.1 .
	docker push go-tg-bot