FROM golang:1.25.1-alpine AS builder

COPY . /go-tg-bot
WORKDIR /go-tg-bot

RUN apk add --no-cache git

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/go-tg-bot-bin cmd/main.go

FROM alpine:3.20.8

WORKDIR /root/
COPY --from=builder /go-tg-bot/bin/go-tg-bot-bin .

CMD ["./go-tg-bot-bin"]