#!/bin/bash
set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"

gcloud config set project $PROJECT_ID
gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

echo -e "\n🔍 Namespaces:"
kubectl get namespaces

echo -e "\n📦 Pods in $CLUSTER_NAME:"
kubectl get pods -n $CLUSTER_NAME -o wide

echo -e "\n🔁 Services:"
kubectl get svc -n $CLUSTER_NAME

echo -e "\n⚙️ Deployments:"
kubectl get deployments -n $CLUSTER_NAME

echo -e "\n📈 Cluster Info:"
kubectl cluster-info
