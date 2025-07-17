# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Clean up unused dependencies and add missing ones
RUN go mod tidy

# Vendor dependencies for reproducible builds
RUN go mod vendor

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o gen-ai-proxy -ldflags "-s -w" ./main.go

# Stage 2: Final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/gen-ai-proxy .

# Copy migration files
COPY --from=builder /app/db/migration ./db/migration

# Copy templates
COPY --from=builder /app/src/templates ./src/templates

# Expose the port your application listens on (e.g., 8080 from your config)
EXPOSE 8080

# Command to run the executable
CMD ["./gen-ai-proxy"]
