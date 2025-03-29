FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/social-service services/social-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/social-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/social-service
RUN go build -o social-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/social-service/social-service .

CMD ["./social-service"]
