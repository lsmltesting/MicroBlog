FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# -------

FROM debian:12-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations /migrations
COPY --from=builder /go/bin/goose /go/bin/goose

ENV PATH="/go/bin:${PATH}"

EXPOSE 8080