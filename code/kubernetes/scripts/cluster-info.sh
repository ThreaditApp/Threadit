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
SHOW_HPA=false
SHOW_DETAILED_RESOURCES=false
SHOW_EVENTS=false

if [[ $# -eq 0 ]]; then
  SHOW_NAMESPACES=true
  SHOW_PODS=true
  SHOW_SERVICES=true
  SHOW_DEPLOYMENTS=true
  SHOW_RESOURCES_PODS=true
  SHOW_RESOURCES_NODES=true
  SHOW_HPA=true
  SHOW_DETAILED_RESOURCES=true
  SHOW_EVENTS=true
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
    --hpa) SHOW_HPA=true ;;
    --detailed-resources) SHOW_DETAILED_RESOURCES=true ;;
    --events) SHOW_EVENTS=true ;;
    --all)
      SHOW_NAMESPACES=true
      SHOW_PODS=true
      SHOW_SERVICES=true
      SHOW_DEPLOYMENTS=true
      SHOW_RESOURCES_PODS=true
      SHOW_RESOURCES_NODES=true
      SHOW_HPA=true
      SHOW_DETAILED_RESOURCES=true
      SHOW_EVENTS=true
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

# New sections for monitoring HPAs and resource limits
$SHOW_HPA && echo -e "\n⚖️ Horizontal Pod Autoscalers:" && kubectl get hpa -n $CLUSTER_NAME

if $SHOW_DETAILED_RESOURCES; then
  echo -e "\n🔎 Detailed Resource Limits and Requests:"
  echo "-------------------------------------------"
  for pod in $(kubectl get pods -n $CLUSTER_NAME -o=name); do
    echo -e "\n📌 $pod"
    kubectl describe $pod -n $CLUSTER_NAME | grep -A8 "Limits:" | grep -v "Node:"
  done
fi

$SHOW_EVENTS && echo -e "\n📜 Recent Events (including scaling):" && kubectl get events -n $CLUSTER_NAME --sort-by='.lastTimestamp' | grep -E '(HorizontalPodAutoscaler|scale|Scaled)'