## 🔍 Phase 5 - System Deployment

The system deployment can be found [here](../../code/kubernetes), divided into:
- [gRPC Gateway](../../code/kubernetes/grpc-gateway) - Contains the deployment and service configuration for the gRPC-Gateway which serves as the RESTful entrypoint to the gRPC-based microservices.
- [MongoDB](../../code/kubernetes/mongo) - Includes manifests to deploy the MongoDB instance used for data persistence by the services.
- [Scripts](../../code/kubernetes/scripts) - A collection of shell scripts used to manage the GKE cluster, deploy services, perform rolling updates and rollbacks, and inspect cluster status.
- [Services](../../code/kubernetes/services) - Includes deployment and service manifests for each of the core backend microservices.
- [Traefik](../../code/kubernetes/traefik) - Configuration files for installing and managing Traefik as an Ingress controller.
- [Config](../../code/kubernetes/config.yaml) - Contains Kubernetes ConfigMap used across services.