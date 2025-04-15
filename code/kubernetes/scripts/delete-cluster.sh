#!/bin/bash
set -e

PROJECT_ID="threadit-api"
CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"

gcloud config set project $PROJECT_ID
gcloud container clusters delete $CLUSTER_NAME --zone=$ZONE --quiet
