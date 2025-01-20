# Use the official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./server ./server
COPY ./main.go ./main.go

# Build the Go app
RUN go build -o app .

# Use a smaller base image for the final container
FROM alpine:latest

# Install Consul
RUN apk add --no-cache curl unzip && \
    curl -o /tmp/consul.zip https://releases.hashicorp.com/consul/1.10.3/consul_1.10.3_linux_amd64.zip && \
    unzip /tmp/consul.zip -d /usr/local/bin/ && \
    rm /tmp/consul.zip

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go app from the builder stage
COPY --from=builder /app/app ./

# Expose ports
EXPOSE 8080 8500

# Add a script to wait for Consul to be ready
COPY wait-for-consul.sh /usr/local/bin/wait-for-consul.sh
RUN chmod +x /usr/local/bin/wait-for-consul.sh

# Command to run the executable
CMD ["sh", "-c", "consul agent -dev -client=0.0.0.0 & wait-for-consul.sh && ./app"]