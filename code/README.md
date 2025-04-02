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

Finally, still in the `/code` directory you need to run the `docker-compose.yml` file. To do that, you should execute the `run.sh` and it will start the necessary containers:

```bash
./run.sh
```

To stop the application execute `stop.sh`:

```bash
./stop.sh
```

## Using gRPC and Protobufs

### Install gRPC and Protobufs for Go

These commands will install all the necessary plugins for Go to utilize protobufs and gRPC, required for Threadit to run correctly 

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### Generate Go Code with Protobufs

Although the generated code is already in this repository, you might want to re-generate it yourself. Verify that you are in the `/code` directory and execute the `generate-proto.sh`, it will generate code for all the required microservices:

```bash
./generate-proto.sh
```

If you wish to generate code for a specific service, you can execute the script with the `-s` flag and give a service name:

```bash
./generate-proto.sh -s service-name
```

### Generate OpenAPI specifications

First use `npm` to install this package:

```bash
npm install -g swagger2openapi 
```

Then just like the generated Go code, follow the same logic and use the script `generate-openapi.sh`:

```bash
./generate-openapi.sh
```

```bash
./generate-openapi.sh -s service-name
```
