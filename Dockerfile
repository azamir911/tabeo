# Stage 1: Build the Go executable
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go executable from the builder stage
COPY --from=builder /app/main /app/main

# Expose the application port
EXPOSE 8080

# Make the executable runnable
RUN chmod +x /app/main

# Run the Go application
CMD ["/app/main"]
