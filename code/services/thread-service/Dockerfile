FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/thread-service services/thread-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/thread-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/thread-service
RUN go build -o thread-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/thread-service/thread-service .

CMD ["./thread-service"]
