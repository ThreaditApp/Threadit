FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/search-service services/search-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/search-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/search-service
RUN go build -o search-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/search-service/search-service .

CMD ["./search-service"]
