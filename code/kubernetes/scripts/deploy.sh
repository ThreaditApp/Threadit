#!/bin/bash
set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"

SKIP_BUILD=false

# Check for --skip-build flag
if [[ "$1" == "--skip-build" ]]; then
  SKIP_BUILD=true
  echo "Skipping image build and push..."
fi

gcloud config set project $PROJECT_ID

# Auth Docker with GCR
gcloud auth configure-docker

# Move to repo code root (Threadit/code/)
cd "$(dirname "$0")/../../"

# Services list
SERVICES=(db-service community-service thread-service comment-service vote-service search-service popular-service)

if [ "$SKIP_BUILD" = false ]; then
  # Build and push all service images
  for SERVICE in "${SERVICES[@]}"; do
    docker build -t gcr.io/$PROJECT_ID/$SERVICE:latest -f services/$SERVICE/Dockerfile .
    docker push gcr.io/$PROJECT_ID/$SERVICE:latest
  done

  # gRPC Gateway
  docker build -t gcr.io/$PROJECT_ID/grpc-gateway:latest -f grpc-gateway/Dockerfile .
  docker push gcr.io/$PROJECT_ID/grpc-gateway:latest
fi

# Move to Kubernetes directory
cd kubernetes

# Authenticate and set up cluster context
gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

# Apply general config
kubectl apply -n $CLUSTER_NAME -f config.yaml

# Traefik
helm repo add traefik https://traefik.github.io/charts
helm upgrade --install traefik traefik/traefik -n $CLUSTER_NAME -f traefik/values.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/cors.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/strip-prefix.yaml

# MongoDB
kubectl apply -n $CLUSTER_NAME -f mongo/

# Services
for SERVICE in "${SERVICES[@]}"; do
  kubectl apply -n $CLUSTER_NAME -f services/$SERVICE/
done

# gRPC Gateway
kubectl apply -n $CLUSTER_NAME -f grpc-gateway/

# Apply Horizontal Pod Autoscalers (HPA)
kubectl apply -n $CLUSTER_NAME -f hpa/
