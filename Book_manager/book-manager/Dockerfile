# Use an official Go runtime as a parent image for building
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the application
COPY . .
RUN go build -o main .

# Use a minimal base image for the runtime
FROM ubuntu:22.04

# Install necessary packages including ca-certificates and wget
RUN apt-get update && \
    apt-get install -y ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

# Download and install dockerize
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz && \
    tar -xzf dockerize-linux-amd64-v0.6.1.tar.gz -C /usr/local/bin && \
    rm dockerize-linux-amd64-v0.6.1.tar.gz

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/views /app/views

# Expose the port on which the application will run
EXPOSE 8080

# Use dockerize to wait for dependencies before starting the application
CMD ["dockerize", "-wait", "tcp://db:5432", "-timeout", "30s", "./main"]
