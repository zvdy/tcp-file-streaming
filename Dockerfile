FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./server ./server

# Build the Go app
RUN go build -o app ./cmd/main.go

# Use a smaller base image for the final container
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go app from the builder stage
COPY --from=builder /app/app ./

# Expose ports 8080 and 8081
EXPOSE 8080
EXPOSE 8081

# Command to run the executable
CMD ["./app"]