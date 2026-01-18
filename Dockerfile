# syntax=docker/dockerfile:1.4

########### Build Stage ###########
FROM golang:1.25 AS builder

WORKDIR /app

# Go mod cache layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Kaynak kodu kopyala
COPY . .

# Statik binary derle
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go


########### Final Stage ###########
FROM alpine:latest

# Gerekli paketler
RUN apk add --no-cache curl

# Non-root user oluştur
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Binary kopyala
COPY --from=builder /app/app .
COPY .env .

# Config dosyalarını kopyala
COPY --from=builder /app/config/ ./config/

# Template dosyalarını kopyala
COPY --from=builder /app/internal/adapters/api/template/ ./internal/adapters/api/template/

# Dosya sahipliğini ayarla
RUN chown -R appuser:appgroup /app

# Non-root kullanıcıya geç
USER appuser

# Çalıştırma komutu
CMD ["./app"]
