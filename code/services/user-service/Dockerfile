FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/user-service services/user-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/user-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/user-service
RUN go build -o user-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/user-service/user-service .

CMD ["./user-service"]
