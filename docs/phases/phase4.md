## 🔍 Phase 4 - System Implementation

The system implementation can be found [here](../../code), divided into:
- [Services](../../code/services/) - The microservices implementation that make up the backend of the system.
- [Protos](../../code/protos/) - The protocol buffers used for gRPC communication between the microservices.
- [Gen](../../code/gen/) - The generated code from the protocol buffers.
- [Traefik](../../code/traefik/) - The configuration for the Traefik reverse proxy used for routing requests to the microservices.
- [gRPC Gateway](../../code/grpc-gateway/) - The configuration for the gRPC Gateway used for exposing the gRPC services as REST APIs.