## Run Threadit on Your Local Machine

### Setup


Threadit requires Go and Docker. You can find installation instructions in the links below:

- [Install Go](https://go.dev/doc/install)

- [Install Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Environmental Variables

To configure the application, change directory to `code/`, then simply copy the example environment file and rename it:

```bash
cp .env.example .env
```

### Run Docker

Finally, still in the `code/` directory you need to run the `docker-compose.yml` file. This command should start the necessary containers:

```bash
docker-compose -p threadit --build -d
```

To stop the application run:

```bash
docker-compose -p threadit stop
```

## Additional information

### Install gRPC and Protobufs for Go

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Generate Go code from the Protobuf Definition

```bash
cd services
./generate-proto.sh <service-name>
```