#!/bin/bash
set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
MACHINE_TYPE="e2-standard-4"
ZONE="europe-west1-b"

gcloud config set project $PROJECT_ID

gcloud container clusters create $CLUSTER_NAME \
  --num-nodes=3 \
  --enable-autoscaling \
  --min-nodes=0 \
  --max-nodes=4 \
  --machine-type=$MACHINE_TYPE \
  --zone=$ZONE \
  --disk-type=pd-standard \
  --disk-size=20

gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

kubectl create ns $CLUSTER_NAME || echo "Namespace already exists."
