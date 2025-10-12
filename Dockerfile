FROM golang:1.23.11-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/ratelimiter cmd/main.go

# FROM scratch
FROM alpine:latest

COPY --from=builder /app .
COPY --from=builder /app/cmd/.env .env

EXPOSE 8080

ENTRYPOINT ["/ratelimiter"]