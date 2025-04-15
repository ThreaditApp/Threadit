#!/bin/bash

# fail on any error
set -e

# check for version argument
if [ -z "$1" ]; then
  echo "❌ Usage: ./update-images.sh <version-tag>"
  exit 1
fi

TAG=$1
PROJECT_ID=$(gcloud config get-value project)
CLUSTER_NAME="threadit-cluster"
NAMESPACE=$CLUSTER_NAME

echo "🔁 Performing rolling updates to version: $TAG"

# update deployments with new image tags
kubectl set image deployment/db-service db-service=gcr.io/$PROJECT_ID/db-service:$TAG -n $NAMESPACE
kubectl set image deployment/community-service community-service=gcr.io/$PROJECT_ID/community-service:$TAG -n $NAMESPACE
kubectl set image deployment/thread-service thread-service=gcr.io/$PROJECT_ID/thread-service:$TAG -n $NAMESPACE
kubectl set image deployment/comment-service comment-service=gcr.io/$PROJECT_ID/comment-service:$TAG -n $NAMESPACE
kubectl set image deployment/vote-service vote-service=gcr.io/$PROJECT_ID/vote-service:$TAG -n $NAMESPACE
kubectl set image deployment/search-service search-service=gcr.io/$PROJECT_ID/search-service:$TAG -n $NAMESPACE
kubectl set image deployment/popular-service popular-service=gcr.io/$PROJECT_ID/popular-service:$TAG -n $NAMESPACE
kubectl set image deployment/grpc-gateway grpc-gateway=gcr.io/$PROJECT_ID/grpc-gateway:$TAG -n $NAMESPACE

echo "✅ Rollout triggered for all services with tag: $TAG"

# show status of rollouts
kubectl rollout status deployment/db-service -n $NAMESPACE
kubectl rollout status deployment/community-service -n $NAMESPACE
kubectl rollout status deployment/thread-service -n $NAMESPACE
kubectl rollout status deployment/comment-service -n $NAMESPACE
kubectl rollout status deployment/vote-service -n $NAMESPACE
kubectl rollout status deployment/search-service -n $NAMESPACE
kubectl rollout status deployment/popular-service -n $NAMESPACE
kubectl rollout status deployment/grpc-gateway -n $NAMESPACE
