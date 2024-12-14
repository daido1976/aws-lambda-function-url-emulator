# Stage 1: Build
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Runtime
FROM alpine:latest

# Set environment variables
ENV RIE_ENDPOINT=http://localhost:9000/2015-03-31/functions/function/invocations

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
