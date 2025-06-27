# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both api and migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate/main.go

# Final image
FROM alpine:latest

WORKDIR /

COPY --from=builder /api /api
COPY --from=builder /migrate /migrate
COPY cmd/migrate/migrations /migrations 

EXPOSE 8080
ENTRYPOINT ["/bin/sh", "-c", "/migrate up && /api"]
