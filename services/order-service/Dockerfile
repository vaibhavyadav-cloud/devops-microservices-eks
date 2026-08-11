# ---------- Stage 1: Build ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Dependency files first for Docker layer caching
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
# Build the application (static binary — no glibc dependency for Alpine runtime)
RUN CGO_ENABLED=0 GOOS=linux go build -o app .
# ---------- Stage 2: Runtime ----------
FROM alpine:3.20
WORKDIR /app
# CA certificates — needed if this service ever makes outbound HTTPS calls
# (not required for Postgres/Kafka today, but cheap to have upfront)
RUN apk add --no-cache ca-certificates
# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# Copy only compiled binary, owned by the non-root user
COPY --from=builder --chown=appuser:appgroup /app/app .
USER appuser
EXPOSE 8082
CMD ["./app"]