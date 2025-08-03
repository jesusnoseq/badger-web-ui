FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o badger-web-ui .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/badger-web-ui .
COPY --from=builder /app/templates ./templates/

ENV BADGER_DB_PATH="/root/badger-data"
ENV BADGER_LOG="false"

EXPOSE 8080
CMD ["./badger-web-ui"]