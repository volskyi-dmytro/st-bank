# Build stage
FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Install migrate
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz -o migrate.tar.gz
RUN tar -xzf migrate.tar.gz
RUN chmod +x migrate

# Run stage
FROM alpine:3.22
WORKDIR /app

# Install netcat for wait-for.sh
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

# Make scripts executable
RUN chmod +x start.sh wait-for.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
