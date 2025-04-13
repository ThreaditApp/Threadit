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
docker build -t gcr.io/$PROJECT_ID/grpc-gateway:latest -f grpc-gateway/Dockerfile .

docker push gcr.io/$PROJECT_ID/db-service:latest
docker push gcr.io/$PROJECT_ID/community-service:latest
docker push gcr.io/$PROJECT_ID/grpc-gateway:latest

CLUSTER_NAME="threadit-cluster"

# create kubernetes cluster
gcloud container clusters create $CLUSTER_NAME \
  --num-nodes=2 \
  --machine-type=e2-medium \
  --zone=europe-west1-b \
  --disk-type=pd-standard \
  --disk-size=10

gcloud container clusters get-credentials $CLUSTER_NAME --zone=europe-west1-b

cd kubernetes

# create a name space
kubectl create ns $CLUSTER_NAME

# Config manifests
kubectl apply -n $CLUSTER_NAME -f config.yaml

# MongoDB manifests
kubectl apply -n $CLUSTER_NAME -f mongo/mongo-pv.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/mongo-secret.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/service.yaml
kubectl apply -n $CLUSTER_NAME -f mongo/deployment.yaml

# DB Service manifests
kubectl apply -n $CLUSTER_NAME -f services/db-service/db-secret.yaml
kubectl apply -n $CLUSTER_NAME -f services/db-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/db-service/deployment.yaml

# Community service manifests
kubectl apply -n $CLUSTER_NAME -f services/community-service/service.yaml
kubectl apply -n $CLUSTER_NAME -f services/community-service/deployment.yaml

# gRPC gateway manifests
kubectl apply -n $CLUSTER_NAME -f grpc-gateway/service.yaml
kubectl apply -n $CLUSTER_NAME -f grpc-gateway/deployment.yaml

# Traefik manifests
kubectl apply -n $CLUSTER_NAME -f traefik/traefik-config.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/service.yaml
kubectl apply -n $CLUSTER_NAME -f traefik/deployment.yaml