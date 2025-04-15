#!/bin/bash
set -e

CLUSTER_NAME="threadit-cluster"
ZONE="europe-west1-b"

# Default: all flags true unless specified
SHOW_NAMESPACES=false
SHOW_PODS=false
SHOW_SERVICES=false
SHOW_DEPLOYMENTS=false
SHOW_RESOURCES_PODS=false
SHOW_RESOURCES_NODES=false

if [[ $# -eq 0 ]]; then
  SHOW_NAMESPACES=true
  SHOW_PODS=true
  SHOW_SERVICES=true
  SHOW_DEPLOYMENTS=true
  SHOW_RESOURCES_PODS=true
  SHOW_RESOURCES_NODES=true
fi

# Parse flags
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --namespaces) SHOW_NAMESPACES=true ;;
    --pods) SHOW_PODS=true ;;
    --services) SHOW_SERVICES=true ;;
    --deployments) SHOW_DEPLOYMENTS=true ;;
    --resources-pods) SHOW_RESOURCES_PODS=true ;;
    --resources-nodes) SHOW_RESOURCES_NODES=true ;;
    --all)
      SHOW_NAMESPACES=true
      SHOW_PODS=true
      SHOW_SERVICES=true
      SHOW_DEPLOYMENTS=true
      SHOW_RESOURCES_PODS=true
      SHOW_RESOURCES_NODES=true
      ;;
    *) echo "❌ Unknown flag: $1"; exit 1 ;;
  esac
  shift
done

# Set project and cluster context
gcloud container clusters get-credentials $CLUSTER_NAME --zone=$ZONE

# Output selected info
$SHOW_NAMESPACES && echo -e "\n🔍 Namespaces:" && kubectl get namespaces
$SHOW_PODS && echo -e "\n📦 Pods in $CLUSTER_NAME:" && kubectl get pods -n $CLUSTER_NAME
$SHOW_SERVICES && echo -e "\n🔁 Services:" && kubectl get svc -n $CLUSTER_NAME
$SHOW_DEPLOYMENTS && echo -e "\n📂 Deployments:" && kubectl get deployments -n $CLUSTER_NAME
$SHOW_RESOURCES_PODS && echo -e "\n📊 Resource Usage (Pods):" && kubectl top pods -n $CLUSTER_NAME
$SHOW_RESOURCES_NODES && echo -e "\n🖥️ Resource Usage (Nodes):" && kubectl top nodes
