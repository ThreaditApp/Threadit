# Kubernetes Deployment Scripts

This folder contains scripts used to deploy the **Threadit API** to **Google Kubernetes Engine (GKE)**.

---

### 1. Create Cluster

Creates a new **GKE cluster** for the Threadit deployment.

```bash
$ ./create-cluster.sh
```
### 2. Deployment 

Builds and pushes Docker images to **Google Container Registry (GCR)**, and deploys all services to the cluster using `kubectl`

```bash
$ ./deploy.sh
```

#### Optional: Skip Image Build & Push

If you've already built and pushed your Docker images, you can skip that step to speed up re-deployments:

```bash
$ ./deploy.sh --skip-build
```

### 3. View Cluster Info 

Displays information about your current GKE cluster and its resources.

```bash
$ ./cluster-info.sh
```

### 4. Delete Cluster

Deletes the Kubernetes cluster and all associated resources to avoid incurring unnecessary costs.

```bash
$ ./delete-cluster.sh
```

---

**Note:** These scripts must be **executed** inside the `scripts/`, otherwise they wil not function properly.