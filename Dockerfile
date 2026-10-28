FROM golang:1.26.4 AS builder
WORKDIR /app

# Install task runner
RUN go install github.com/go-task/task/v3/cmd/task@b5aef1c287a9be8038347a17646bdf4bc7185bce

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY ./ .

# Ensure vendor directory is up-to-date after copying source
RUN go mod vendor

# Install migrator
RUN task install:migrator

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/backend cmd/backend/main.go

# Production stage
FROM debian:bookworm-slim
WORKDIR /app

# Install only essential packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    wget && \
    rm -rf /var/lib/apt/lists/*

# Copy necessary files from builder
COPY --from=builder /app/taskfile*.yml ./
COPY --from=builder /app/bin/backend ./bin/backend
COPY --from=builder /go/bin/task /usr/local/bin/task
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY db/migrations ./db/migrations

# Create non-root user for security
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser
