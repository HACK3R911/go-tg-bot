FROM golang:1.25.6-alpine AS builder

COPY . /go-tg-bot
WORKDIR /go-tg-bot

RUN apk add --no-cache git

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/go-tg-bot-bin ./cmd/bot

FROM alpine:3.20

WORKDIR /root/
COPY --from=builder /go-tg-bot/bin/go-tg-bot-bin .
COPY --from=builder /go-tg-bot/migrations ./migrations

# Note: Configuration should be provided via environment variables
# Do NOT copy .env file - use docker-compose env_file or environment variables

CMD ["./go-tg-bot-bin"]