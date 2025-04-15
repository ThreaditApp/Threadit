#!/bin/bash

# TODO this script is not yet ready, needs more testing

set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"
NAMESPACE=$CLUSTER_NAME

SERVICES=(db-service community-service thread-service comment-service vote-service search-service popular-service grpc-gateway)

print_usage() {
  echo "Usage:"
  echo "  $0 <service-name> <image-tag>"
  echo "  $0 --all <image-tag>"
  echo ""
  echo "Examples:"
  echo "  $0 community-service latest"
  echo "  $0 --all latest"
}

if [[ $# -lt 2 ]]; then
  echo "❌ Not enough arguments."
  print_usage
  exit 1
fi

if [[ "$1" == "--all" ]]; then
  TAG="$2"
  for SERVICE in "${SERVICES[@]}"; do
    echo -e "\n🚀 Performing rolling update for '$SERVICE' to image 'gcr.io/$PROJECT_ID/$SERVICE:$TAG'..."
    kubectl set image deployment/$SERVICE $SERVICE=gcr.io/$PROJECT_ID/$SERVICE:$TAG -n $NAMESPACE
    echo "⏳ Waiting for '$SERVICE' rollout to complete..."
    kubectl rollout status deployment/$SERVICE -n $NAMESPACE
  done
  echo -e "\n✅ All deployments updated to tag '$TAG'"
else
  SERVICE="$1"
  TAG="$2"
  echo "🚀 Performing rolling update for deployment '$SERVICE' with image 'gcr.io/$PROJECT_ID/$SERVICE:$TAG'..."
  kubectl set image deployment/$SERVICE $SERVICE=gcr.io/$PROJECT_ID/$SERVICE:$TAG -n $NAMESPACE
  echo "⏳ Waiting for rollout to complete..."
  kubectl rollout status deployment/$SERVICE -n $NAMESPACE
  echo "✅ Rolling update complete!"
fi
