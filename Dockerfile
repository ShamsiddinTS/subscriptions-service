FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscriptions-service ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/subscriptions-service .
COPY --from=builder /app/internal/migrations ./internal/migrations
COPY --from=builder /app/internal/config ./internal/config
COPY .env .

EXPOSE 8080

CMD ["./subscriptions-service"]