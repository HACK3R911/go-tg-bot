build:
	GOOS=linux GOARCH=amd64 go build -o ./bin cmd/bot/main.go

docker-build-n-push:
	docker buildx build --no-cache --platform linux/amd64 -t hack3r11/go-tg-bot:v0.0.1 .
	docker login -u hack3r11 --password-stdin $DOCKER_HUB_TOKEN
	docker push hack3r11/go-tg-bot:v0.0.1