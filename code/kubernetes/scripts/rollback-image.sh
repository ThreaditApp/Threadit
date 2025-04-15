#!/bin/bash

# TODO this script is not yet ready, needs more testing

set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"
NAMESPACE=$CLUSTER_NAME

# List of all services
SERVICES=(db-service community-service thread-service comment-service vote-service search-service popular-service grpc-gateway)

print_usage() {
  echo "Usage:"
  echo "  $0 <service-name>"
  echo "  $0 --all"
  echo ""
  echo "Examples:"
  echo "  $0 community-service"
  echo "  $0 --all"
}

# Validate arguments
if [[ $# -lt 1 ]]; then
  echo "❌ Not enough arguments."
  print_usage
  exit 1
fi

# Rollback logic
if [[ "$1" == "--all" ]]; then
  for SERVICE in "${SERVICES[@]}"; do
    echo -e "\n⏪ Rolling back deployment '$SERVICE'..."
    kubectl rollout undo deployment/$SERVICE -n $NAMESPACE
    echo "⏳ Waiting for '$SERVICE' to roll out previous revision..."
    kubectl rollout status deployment/$SERVICE -n $NAMESPACE
  done
  echo -e "\n✅ All deployments rolled back to previous revisions."
else
  SERVICE="$1"
  echo "⏪ Rolling back deployment '$SERVICE'..."
  kubectl rollout undo deployment/$SERVICE -n $NAMESPACE
  echo "⏳ Waiting for rollout to complete..."
  kubectl rollout status deployment/$SERVICE -n $NAMESPACE
  echo "✅ Rollback complete for '$SERVICE'."
fi