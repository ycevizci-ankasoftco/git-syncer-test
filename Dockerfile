# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o git-syncer

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/git-syncer .

ENTRYPOINT ["./git-syncer"]
