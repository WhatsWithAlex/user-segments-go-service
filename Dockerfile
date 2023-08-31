# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Run stage
FROM scratch
COPY --from=builder /bin/app /app
COPY --from=builder /app/.env /.env
COPY --from=builder /app/sql/migrations /sql/migrations
COPY --from=builder /app/web/static /web/static

EXPOSE ${APP_PORT}
CMD [ "/app" ]