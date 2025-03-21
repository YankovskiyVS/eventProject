# Stage 1: Build the application
FROM golang:alpine AS builder

# Set the working directory
WORKDIR /app

# Copy module files first for layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project
COPY . .

# Build from cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Stage 2: Runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8000
CMD ["./main"]