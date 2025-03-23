## Run Threadit on Your Local Machine

### Setup


Threadit requires Go and Docker. You can find installation instructions in the links below:

- [Install Go](https://go.dev/doc/install)

- [Install Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Environmental Variables

To configure the application, change directory to `code/`, create a `.env` file in the root of this directory. The `.env` file should contain necessary environment variables for running the application.

```ini
# Example .env configuration
MONGO_INITDB_DATABASE=threadit-db
MONGO_INITDB_ROOT_USERNAME=test
MONGO_INITDB_ROOT_PASSWORD=example
MONGO_URI=mongodb://test:example@mongodb:27017/threadit-db?authSource=test
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