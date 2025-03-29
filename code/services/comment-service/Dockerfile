FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/comment-service services/comment-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/comment-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/comment-service
RUN go build -o comment-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/comment-service/comment-service .

CMD ["./comment-service"]
