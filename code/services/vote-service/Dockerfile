FROM golang:1.23 AS builder

WORKDIR /app

# Copy the folder with the generated code
COPY gen/ gen/

# Copy the folder with the service source code
COPY services/vote-service services/vote-service

# Download dependencies
WORKDIR /app/gen
RUN go mod download
WORKDIR /app/services/vote-service
RUN go mod download

# Build the service binary
WORKDIR /app/services/vote-service
RUN go build -o vote-service .

# Use a minimal runtime environment
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/services/vote-service/vote-service .

CMD ["./vote-service"]
