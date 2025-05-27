#!/bin/bash
set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
MACHINE_TYPE="e2-standard-2"
ZONE="europe-west1-b"

gcloud config set project $PROJECT_ID

gcloud container clusters create $CLUSTER_NAME \
  --num-nodes=3 \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=5 \
  --machine-type=$MACHINE_TYPE \
  --zone=$ZONE \
  --disk-type=pd-standard \
  --disk-size=25

gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

kubectl create ns $CLUSTER_NAME || echo "Namespace already exists."
