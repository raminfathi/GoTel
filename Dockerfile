# Stage 1: Build the application
# استفاده از نسخه 1.25 (یا latest)
FROM golang:1.25-alpine as builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# We explicitly disable CGO for a static binary and point to main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Stage 2: Run the application
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 5000 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["./main"]