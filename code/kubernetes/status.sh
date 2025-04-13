#!/bin/bash

# fail on any error
set -e

CLUSTER_NAME="threadit-cluster"

# Show all resources in the namespace
kubectl get pods -n $CLUSTER_NAME

# Grab external ip
kubectl get service traefik -n $CLUSTER_NAME