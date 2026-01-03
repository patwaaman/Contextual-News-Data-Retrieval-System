# ===== Build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o news-retrieval main.go

# ===== Runtime stage =====
FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/news-retrieval /app/news-retrieval

EXPOSE 8080
EXPOSE 9090

CMD ["./news-retrieval"]
