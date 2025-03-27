FROM golang:1.24.1 AS builder

WORKDIR /app

# Copy and download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the service binary
RUN go build -o thread-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/thread-service .

CMD ["./thread-service"]
