# Build stage
FROM golang:1.23.4-alpine AS builder
# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest  

# Set the working directory
WORKDIR /root/

# Create logs directory
RUN mkdir -p /root/logs && \
    chmod 755 /root/logs

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the .env file
COPY .env .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]