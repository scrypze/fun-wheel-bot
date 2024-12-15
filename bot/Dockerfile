FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init fun-wheel-bot && go mod tidy && go build -o bot .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
ENV TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
CMD ["./bot"]
