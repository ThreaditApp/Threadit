## Run Threadit on Your Local Machine

### Setup

Threadit requires Go and Docker. You can find installation instructions in the links below:

- [Install Go](https://go.dev/doc/install)
- [Install Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Pulling Datasets

The datasets in this repository are tracked using Git LFS (Large File Storage). To download them, you need to run:

```bash
git lfs install
git lfs pull
```

### Environment Variables

To configure the application, change directory to `/code`, then simply copy the example environment file and rename it:

```bash
cp .env.example .env
```

### Run Docker

Finally, still in the `/code` directory you need to run the `docker-compose.yml` file. This command should start the necessary containers:

```bash
docker-compose -p threadit up --build -d
```

To stop the application run:

```bash
docker-compose -p threadit stop
```

## Using gRPC and Protobufs

### Install gRPC and Protobufs for Go

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### Generate Go Code for a Service with Protobufs

```bash
cd code
./generate-proto.sh service-name
```

### Generate Go Code for All Services with Protobufs
```bash
cd code
./generate-proto-all.sh
```
