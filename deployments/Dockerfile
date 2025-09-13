# Multi-stage Dockerfile for Domain MAX

# Stage 1: Build frontend
FROM node:18-alpine AS web-builder

# Set working directory
WORKDIR /app/web

# Copy package files for better caching
COPY web/package*.json ./

# Install dependencies with npm ci for production builds
RUN npm ci --only=production --no-audit --no-fund

# Copy frontend source code
COPY web/ ./

# Build frontend application
RUN npm run build

# Verify build output
RUN ls -la dist/ && echo "Frontend build completed successfully"

# Stage 2: Build backend
FROM golang:1.23-alpine AS server-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Set working directory
WORKDIR /app

# Copy Go module files for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

# Copy built frontend files from previous stage
COPY --from=web-builder /app/web/dist ./web/dist

# Verify frontend files are copied
RUN ls -la web/dist/ && echo "Frontend files copied successfully"

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o domain-max ./cmd/server

# Verify binary
RUN file domain-max && ls -lh domain-max

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    wget \
    curl \
    tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/*

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=server-builder /app/domain-max .
COPY --from=server-builder /app/web/dist ./web/dist

# Copy configuration files if they exist
COPY --from=server-builder /app/configs ./configs

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Add metadata
LABEL maintainer="Domain MAX Team" \
      version="1.0.0" \
      description="Domain MAX - DNS Management System"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

<<<<<<< HEAD
# Start the application
CMD ["./domain-max"]
=======
# 启动应用
CMD ["./domain-max"]
>>>>>>> 2c9f9981e8704ab02eee3aa85091a742a44e67f9
