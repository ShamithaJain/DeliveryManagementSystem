# Use official Go image as the builder
FROM golang:1.21-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set work directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o dms ./cmd

# Final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/dms .

# Copy .env file if exists
COPY .env .

# Expose port
EXPOSE 8080

# Command to run the app
CMD ["./dms"]
