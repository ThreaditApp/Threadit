#!/bin/bash

# fail on any error
set -e

# set gcloud project
gcloud config set project threadit-api
PROJECT_ID=$(gcloud config get-value project)
gcloud config get-value project

cd ..

# build and push container images
gcloud auth configure-docker

docker build -t gcr.io/$PROJECT_ID/db-service:latest -f services/db-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/community-service:latest -f services/community-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/thread-service:latest -f services/thread-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/comment-service:latest -f services/comment-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/vote-service:latest -f services/vote-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/search-service:latest -f services/search-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/popular-service:latest -f services/popular-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/grpc-gateway:latest -f grpc-gateway/Dockerfile .

docker push gcr.io/$PROJECT_ID/db-service:latest
docker push gcr.io/$PROJECT_ID/community-service:latest
docker push gcr.io/$PROJECT_ID/thread-service:latest
docker push gcr.io/$PROJECT_ID/comment-service:latest
docker push gcr.io/$PROJECT_ID/vote-service:latest
docker push gcr.io/$PROJECT_ID/search-service:latest
docker push gcr.io/$PROJECT_ID/popular-service:latest
docker push gcr.io/$PROJECT_ID/grpc-gateway:latest

# create kubernetes cluster
CLUSTER_NAME="threadit-cluster"
gcloud container clusters create $CLUSTER_NAME \
  --num-nodes=2 \
  --machine-type=e2-standard-2 \
  --zone=europe-west1-b \
  --disk-type=pd-standard \
  --disk-size=10

gcloud container clusters get-credentials $CLUSTER_NAME --zone=europe-west1-b

# create a name space
kubectl create ns $CLUSTER_NAME

# Apply kubernetes manifests
cd kubernetes

# Config
kubectl apply -n $CLUSTER_NAME -f config.yaml

# MongoDB
kubectl apply -n $CLUSTER_NAME -f mongo/mongo-pv.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/mongo-secret.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/service.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/deployment.yaml

# DB service
kubectl apply -n $CLUSTER_NAME -f services/db-service/db-secret.yaml
kubectl apply -n $CLUSTER_NAME -f services/db-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/db-service/deployment.yaml

# Community service
kubectl apply -n $CLUSTER_NAME -f services/community-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/community-service/deployment.yaml

# Thread service
kubectl apply -n $CLUSTER_NAME -f services/thread-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/thread-service/deployment.yaml

# Comment service
kubectl apply -n $CLUSTER_NAME -f services/comment-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/comment-service/deployment.yaml

# Vote service
kubectl apply -n $CLUSTER_NAME -f services/vote-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/vote-service/deployment.yaml

# Search service
kubectl apply -n $CLUSTER_NAME -f services/search-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/search-service/deployment.yaml

# Popular service
kubectl apply -n $CLUSTER_NAME -f services/popular-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/popular-service/deployment.yaml

# gRPC gateway
kubectl apply -n $CLUSTER_NAME -f grpc-gateway/service.yaml
kubectl apply -n $CLUSTER_NAME -f grpc-gateway/deployment.yaml

# Traefik
kubectl apply -n $CLUSTER_NAME -f traefik/traefik-config.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/service.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/deployment.yaml