#!/bin/bash

# fail on any error
set -e

CLUSTER_NAME="threadit-cluster"
NAMESPACE=$CLUSTER_NAME

echo "⏪ Rolling back all services in namespace: $NAMESPACE"

kubectl rollout undo deployment/db-service -n $NAMESPACE
kubectl rollout undo deployment/community-service -n $NAMESPACE
kubectl rollout undo deployment/thread-service -n $NAMESPACE
kubectl rollout undo deployment/comment-service -n $NAMESPACE
kubectl rollout undo deployment/vote-service -n $NAMESPACE
kubectl rollout undo deployment/search-service -n $NAMESPACE
kubectl rollout undo deployment/popular-service -n $NAMESPACE
kubectl rollout undo deployment/grpc-gateway -n $NAMESPACE

echo "✅ Rollbacks triggered"

# show status of rollouts
kubectl rollout status deployment/db-service -n $NAMESPACE
kubectl rollout status deployment/community-service -n $NAMESPACE
kubectl rollout status deployment/thread-service -n $NAMESPACE
kubectl rollout status deployment/comment-service -n $NAMESPACE
kubectl rollout status deployment/vote-service -n $NAMESPACE
kubectl rollout status deployment/search-service -n $NAMESPACE
kubectl rollout status deployment/popular-service -n $NAMESPACE
kubectl rollout status deployment/grpc-gateway -n $NAMESPACE
