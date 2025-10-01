# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the httpserver binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/httpserver ./cmd/httpserver

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/httpserver .

# Expose the application port
EXPOSE 42069

# Run the application
CMD ["./httpserver"] 