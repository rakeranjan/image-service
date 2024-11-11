FROM golang:1.22.8 AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the application binary
WORKDIR /app/cmd
RUN go build -o /app/backend

# Stage 2: Run
FROM debian:bookworm-slim

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/backend /app/backend

# Expose the application's port
EXPOSE 8001

# Command to run the application
CMD ["./backend"]