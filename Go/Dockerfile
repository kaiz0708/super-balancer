# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o load-balancer main.go

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/load-balancer .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./load-balancer"]
