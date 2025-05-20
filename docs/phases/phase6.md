## 🔍 Phase 6 - Non-Functional Requirements and Technical Architecture or Big Data

## 1. Planned Improvements

### 1.1 Ingress/API Gateway Configuration

Currently, we are using Traefik with a basic IngressRoute and minimal configuration. To improve flexibility and take advantage of Kubernetes-native features, we will explore two alternatives:

- **Kubernetes Ingress with Traefik:** This allows using standard Ingress resources with Traefik's CRDs for fine-grained traffic routing, TLS termination, and middleware chaining.

- **Kubernetes Gateway API:** A more expressive and extensible alternative to the Ingress API, which decouples traffic routing from infrastructure. We'll experiment with it alongside Traefik's Gateway API support to modernize our networking layer.

### 1.2 Liveness and Readiness Probes
To improve fault tolerance and enable better self-healing behavior in Kubernetes, we will define:

* **Liveness probes** to detect and restart failed pods.
* **Readiness probes** to ensure that traffic is only sent to pods that are ready to handle requests.

### 1.3 Resource Limits and HPA
We will benchmark services to determine ideal values for:

* **CPU and memory resource requests/limits**.
* **Horizontal Pod Autoscaling (HPA)** thresholds based on real traffic patterns to ensure scalability.

### 1.4 Authentication & Authorization with Keycloak
We will introduce authentication and authorization by integrating [Keycloak](https://www.keycloak.org/) as the Identity Provider. Keycloak will manage user sessions, tokens (OIDC) and Role-Based Access Control (RBAC) across the services.

### 1.5 Secret management

To further improve security we will explore Google Secret Manager for managing sensitive configuration data such as API keys, credentials and tokens. This approach provides:

- Centralized and secure secret storage.

- IAM-based access control.

- Seamless integration with GKE via workload identity.

- Versioning and audit logging for secret access.

### 1.6 CI/CD Pipeline
A continuous integration and deployment (CI/CD) pipeline will be implemented using GitHub Actions and GKE in order to automatically:

* Update OpenAPI specifications.
* Build images.
* Deploy services to the cluster.

## 2. Non-Functional Requirements
| Category        | Requirement                                                                |
| --------------- | -------------------------------------------------------------------------- |
| Scalability     | Use HPA to autoscale services based on CPU usage.                          |
| Availability    | Configure liveness/readiness probes and multiple replicas where necessary. |
| Security        | Enforce authentication and authorization via Keycloak.                     |
| Maintainability | Implement CI/CD for consistent, automated deployments.                     |

## 3. Deployment Plan
| Step                            | Tool/Technology       |
| ------------------------------- | --------------------- |
| Containerization                | Docker                |
| Cluster Orchestration           | Kubernetes (GKE)      |
| Ingress Management              | GKE Native Ingress    |
| Identity and Access Management  | Keycloak (OIDC, RBAC) |
| CI/CD                           | GitHub Actions + GKE  |
| Autoscaling                     | Kubernetes HPA        |

## 4. Architecture Diagram
![application architecture](../images/architecture.png)
