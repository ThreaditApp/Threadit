#!/bin/bash
set -e

BUILD=false
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"
SERVICES=(db community thread comment vote search popular)

# Set project and set up cluster context
gcloud config set project $PROJECT_ID
gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

GCS_KEY="gcs-key"
BUCKET_SECRET=$(gcloud secrets versions access latest --secret=$GCS_KEY)
MONGO_USER=$(gcloud secrets versions access latest --secret="mongo-user")
MONGO_PASS=$(gcloud secrets versions access latest --secret="mongo-pass")

# Check for --build flag
if [[ "$1" == "--build" ]]; then
  BUILD=true
  echo "Building and pushing images..."
fi

# Build and push docker images
build_and_push_images() {
  cd "$SCRIPT_DIR/../../" || exit 1

  gcloud auth configure-docker

  for SERVICE in "${SERVICES[@]}"; do
    docker build -t gcr.io/$PROJECT_ID/"$SERVICE-service":latest -f services/"$SERVICE-service"/Dockerfile .
    docker push gcr.io/$PROJECT_ID/"$SERVICE-service":latest
  done

  docker build -t gcr.io/$PROJECT_ID/grpc-gateway:latest -f grpc-gateway/Dockerfile .
  docker push gcr.io/$PROJECT_ID/grpc-gateway:latest

  cd "$SCRIPT_DIR" || exit 1
}

# Build and push images if --build is passed
if [ "$BUILD" = true ]; then
  build_and_push_images
fi

cd "$SCRIPT_DIR/.." || exit 1

# Deploy traefik
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik -n $CLUSTER_NAME -f traefik/values.yaml

kubectl apply -n $CLUSTER_NAME -f traefik/cors.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/strip-prefix.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/ingress-routes.yaml

# Deploy threadit application
kubectl create secret generic "bucket-secret" \
  --from-literal="$GCS_KEY.json=$BUCKET_SECRET" \
  -n $CLUSTER_NAME --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic "mongo-secret" \
  --from-literal="MONGO_INITDB_ROOT_USERNAME=$MONGO_USER" \
  --from-literal="MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASS" \
  -n $CLUSTER_NAME --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -n $CLUSTER_NAME -f config.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/

echo "Deploying Keycloak..."
kubectl apply -n $CLUSTER_NAME -f keycloak/configmap.yaml
kubectl apply -n $CLUSTER_NAME -f keycloak/secrets.yaml
kubectl apply -n $CLUSTER_NAME -f keycloak/realm-configmap.yaml
kubectl apply -n $CLUSTER_NAME -f keycloak/deployment.yaml

for SERVICE in "${SERVICES[@]}"; do
  kubectl apply -n $CLUSTER_NAME -f services/"$SERVICE-service"/
done

kubectl apply -n $CLUSTER_NAME -f grpc-gateway/

echo "Waiting for Keycloak to be ready..."
kubectl wait --for=condition=ready pod -l app=keycloak -n $CLUSTER_NAME --timeout=300s

echo "Deployment complete!"