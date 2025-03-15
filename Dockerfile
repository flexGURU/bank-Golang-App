# Build stage
FROM golang:1.23.6-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migration


EXPOSE 8000
CMD [ "/app/main" ]

ENTRYPOINT [ "/app/start.sh" ]