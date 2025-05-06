## 🔍 Phase 6 - Non-Functional Requirements and Technical Architecture or Big Data

## 1. Planned Improvements

### 1.1 Architecture Update
We plan to update the architecture diagram to provide a more accurate representation of:

* The internal communication between microservices via gRPC.
* The gRPC Gateway responsible for converting REST requests to gRPC.
* The ingress flow, now using GKE’s native Ingress controller instead of Traefik.

### 1.2 Migration to GKE Ingress
To simplify deployment and maintenance, we will remove the Traefik Ingress Controller and adopt the native GKE Ingress. This improves integration with GCP’s load balancing.

### 1.3 Authentication with Keycloak
We will introduce authentication and authorization by integrating [Keycloak](https://www.keycloak.org/) as the Identity Provider. Keycloak will manage user sessions, tokens (OIDC), and role-based access control (RBAC) across the services.

### 1.4 CI/CD Pipeline
A continuous integration and deployment (CI/CD) pipeline will be implemented using GitHub Actions and GKE. It will:

* Update OpenAPI specification files.
* Build images.
* Deploy services to the cluster.

### 1.5 Liveness and Readiness Probes
To improve fault tolerance and enable better self-healing behavior in Kubernetes, we will define:

* **Liveness probes** to detect and restart failed containers.
* **Readiness probes** to ensure that traffic is only sent to containers that are ready to handle requests.

### 1.6 Resource Limits and HPA
We will benchmark services to determine ideal values for:

* **CPU and memory resource requests/limits**.
* **Horizontal Pod Autoscaling (HPA)** thresholds to ensure scalability based on real traffic patterns.

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