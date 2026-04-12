# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Enable go modules
ENV GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o food-delivery ./cmd/server

# Stage 2: Run minimal image
FROM alpine:latest  

WORKDIR /app/

# Copy binary from builder
COPY --from=builder /app/food-delivery .

EXPOSE 8080

CMD ["./food-delivery"]
