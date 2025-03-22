
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