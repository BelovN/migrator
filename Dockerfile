FROM golang:1.23-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o migrator .

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y libc6 ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/migrator /usr/local/bin/migrator

ENTRYPOINT ["/usr/local/bin/migrator"]