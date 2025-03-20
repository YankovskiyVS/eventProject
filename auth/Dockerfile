# Stage 1: Build the application
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy environment variables (if needed)
COPY .env .

# Expose the port the application will run on
EXPOSE 8000

# Command to run the application
CMD ["./main"]