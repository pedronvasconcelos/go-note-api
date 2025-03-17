# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install swag and generate Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go -o internal/api/docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o note-api ./cmd/api/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/note-api .
RUN ls -la ./note-api
# Set executable permissions
RUN chmod +x ./note-api

# Expose port
EXPOSE 8080

# Command to run
CMD ["./note-api"] 