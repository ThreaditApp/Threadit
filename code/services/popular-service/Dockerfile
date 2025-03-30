FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/popular-service services/popular-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/popular-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/popular-service
RUN go build -o popular-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/popular-service/popular-service .

CMD ["./popular-service"]
