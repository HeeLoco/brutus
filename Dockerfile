# Build stage
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy go modules and build the application
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o brutus cmd/brutus/main.go

# Production stage
FROM ubuntu:24.04

# Install Azure CLI and basic tools
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    jq \
    unzip \
    && curl -sL https://aka.ms/InstallAzureCLIDeb | bash \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -s /bin/bash brutus

# Copy the binary from builder stage
COPY --from=builder /app/brutus /usr/local/bin/brutus

# Switch to non-root user
USER brutus
WORKDIR /home/brutus

# Set environment variables
ENV TERM=xterm-256color

# Default command
CMD ["brutus"]