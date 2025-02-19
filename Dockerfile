# Use the official Golang image as a build stage
FROM golang:1.23.6-alpine3.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Set the Go proxy to speed up module download
ENV GOPROXY=https://goproxy.cn,direct

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final stage
FROM alpine:3.21.3

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application will run on
EXPOSE 8000

# Command to run the application
CMD ["./main", "serve"] 