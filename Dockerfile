# Multi-stage Dockerfile for Cloud Run (zero CGO / pure-Go SQLite)
FROM golang:alpine AS builder

WORKDIR /src

# Install CA certificates for external HTTPS calls if needed
RUN apk add --no-cache ca-certificates git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build pure-Go statically linked binary (no CGO required)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /bin/precedent ./cmd/server

# Final minimal deployment image
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy binary and static templates
COPY --from=builder /bin/precedent /app/precedent
COPY --from=builder /src/web/templates /app/web/templates

# Default environment for Cloud Run
ENV PORT=8080
ENV TEMPLATE_DIR=/app/web/templates
ENV DB_PATH=/app/precedent.db

EXPOSE 8080

ENTRYPOINT ["/app/precedent"]
