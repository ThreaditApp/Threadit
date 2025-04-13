#!/bin/bash

# fail on any error
set -e

CLUSTER_NAME="threadit-cluster"
gcloud container clusters delete ${CLUSTER_NAME} --zone=europe-west1-b --quiet
