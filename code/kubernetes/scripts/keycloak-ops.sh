#!/bin/bash
set -e

CLUSTER_NAME="threadit-cluster"

function help() {
    echo "Usage: $0 <command>"
    echo "Commands:"
    echo "  status       - Check Keycloak status"
    echo "  logs         - Show Keycloak logs"
    echo "  restart      - Restart Keycloak deployment"
    echo "  reload       - Reload realm configuration"
    echo "  port-forward - Start port forwarding to access Keycloak locally"
}

function check_status() {
    echo "Checking Keycloak status..."
    kubectl get pods -n $CLUSTER_NAME -l app=keycloak
}

function show_logs() {
    echo "Fetching Keycloak logs..."
    kubectl logs -n $CLUSTER_NAME -l app=keycloak --tail=100 -f
}

function restart_keycloak() {
    echo "Restarting Keycloak..."
    kubectl rollout restart deployment/keycloak -n $CLUSTER_NAME
    kubectl rollout status deployment/keycloak -n $CLUSTER_NAME
}

function reload_realm() {
    echo "Reloading realm configuration..."
    # Delete the existing pod to force a reload of the realm config
    kubectl delete pod -n $CLUSTER_NAME -l app=keycloak
    echo "Waiting for new pod to be ready..."
    kubectl wait --for=condition=ready pod -l app=keycloak -n $CLUSTER_NAME --timeout=300s
}

function port_forward() {
    echo "Starting port forward to Keycloak on localhost:8080..."
    kubectl port-forward -n $CLUSTER_NAME svc/keycloak 8080:8080
}

case "$1" in
    "status")
        check_status
        ;;
    "logs")
        show_logs
        ;;
    "restart")
        restart_keycloak
        ;;
    "reload")
        reload_realm
        ;;
    "port-forward")
        port_forward
        ;;
    *)
        help
        exit 1
        ;;
esac 