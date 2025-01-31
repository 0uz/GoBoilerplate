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

# Add necessary packages
RUN apk add --no-cache curl

# Create a non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the .env file
COPY .env .

# Set ownership of the application files
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Command to run the executable
CMD ["./main"]