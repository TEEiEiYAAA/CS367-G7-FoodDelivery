# Stage 1: Build the Go binary
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o food-delivery ./cmd/server

# Stage 2: Runtime image
FROM alpine:latest

WORKDIR /app/

# Copy binary from builder
COPY --from=builder /app/food-delivery .

EXPOSE 8080

CMD ["./food-delivery"]
