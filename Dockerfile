# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code
COPY main.go ./

# Build the application
RUN go build -o loadtest main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/loadtest .

# Add execute permissions
RUN chmod +x loadtest

# Set the entrypoint
ENTRYPOINT ["./loadtest"]