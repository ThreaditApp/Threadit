FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY grpc-gateway grpc-gateway

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/grpc-gateway
RUN go mod download

# Build the service binary
WORKDIR /app/grpc-gateway
RUN go build -o grpc-gateway .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/grpc-gateway/grpc-gateway .

CMD ["./grpc-gateway"]
