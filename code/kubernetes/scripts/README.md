# Kubernetes Deployment Scripts

This folder contains scripts used to deploy the **Threadit API** to **Google Kubernetes Engine (GKE)**.

---

### 1. Create Cluster

Creates a new **GKE cluster** for the Threadit deployment.

```bash
$ ./create-cluster.sh
```
### 2. Deployment 

Deploys all services and components to the cluster using `kubectl`.

```bash
$ ./deploy.sh
```

This will deploy using the latest available Docker images in **Google Container Registry (GCR)**.

#### Optional: Build & Push

If you want to build and push the Docker images before deploying, use the --build flag:

```bash
$ ./deploy.sh --build
```

### 3. View Cluster Info

Displays information about your current GKE cluster and its resources.

```bash
$ ./cluster-info.sh
```

#### Options:

- `--namescapes` Displays a list of all Kubernetes namespaces in the current cluster.
- `--pods` Shows the status and details of all pods running in the specified namespace.
- `--services` Lists all services deployed in the specified namespace.
- `--deployments` Shows deployment configurations and statuses for the namespace.
- `--resources-pods` Displays real-time CPU and memory usage metrics for each pod.
- `--resources-nodes` Displays real-time CPU and memory usage metrics for each node in the cluster.
- `--all` Runs all of the above commands to display full cluster info.

### 4. Delete Cluster

Deletes the Kubernetes cluster and all associated resources to avoid incurring unnecessary costs.

```bash
$ ./delete-cluster.sh
```

---

**Note:** These scripts must be **executed** inside the `scripts/`, otherwise they wil not function properly.