#!/bin/bash

# fail on any error
set -e

# check for a service name argument
if [ -z "$1" ]; then
  echo "❌ Usage: ./generate-proto.sh <service-name>"
  exit 1
fi

SERVICE_NAME="$1"

# get paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="$SCRIPT_DIR/proto"
SERVICE_DIR="$SCRIPT_DIR/services/$SERVICE_NAME"
PROTO_FILE="$PROTO_DIR/${SERVICE_NAME}.proto"
OUT_DIR="$SERVICE_DIR/src/pb"
REQUIREMENTS_FILE="$SERVICE_DIR/requirements.txt"

# check if the .proto file exists
if [ ! -f "$PROTO_FILE" ]; then
  echo "❌ Proto file not found: $PROTO_FILE"
  exit 1
fi

echo "🔄 Generating Go code from $PROTO_FILE..."

# create output directory if it doesn't exist
mkdir -p "$OUT_DIR"

# clean up existing generated files
rm -f "$OUT_DIR"/*.pb.go

# run protoc
protoc \
  --go_out="$OUT_DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" \
  --go-grpc_opt=paths=source_relative \
  --proto_path="$PROTO_DIR" \
  "$PROTO_FILE"

# process dependencies
if [ -f "$REQUIREMENTS_FILE" ]; then
  echo "📋 Found requirements.txt, processing dependencies..."
  
  if command -v dos2unix &> /dev/null; then
    dos2unix "$REQUIREMENTS_FILE" &> /dev/null || true
  fi
  
  while IFS= read -r DEPENDENCY || [[ -n "$DEPENDENCY" ]]; do
    if [[ -z "$DEPENDENCY" || "$DEPENDENCY" =~ ^# || "$DEPENDENCY" =~ ^// ]]; then
      continue
    fi
    
    DEPENDENCY=$(echo "$DEPENDENCY" | xargs)
    
    echo "🔄 Processing dependency: $DEPENDENCY"
    DEP_PROTO_FILE="$PROTO_DIR/${DEPENDENCY}.proto"
    
    if [ ! -f "$DEP_PROTO_FILE" ]; then
      echo "⚠️  Warning: Proto file not found for dependency: $DEP_PROTO_FILE"
      continue
    fi
    
    protoc \
      --go_out="$OUT_DIR" \
      --go_opt=paths=source_relative \
      --go-grpc_out="$OUT_DIR" \
      --go-grpc_opt=paths=source_relative \
      --proto_path="$PROTO_DIR" \
      "$DEP_PROTO_FILE"
      
    echo "✅ Generated Go code for dependency '$DEPENDENCY'"
  done < "$REQUIREMENTS_FILE"
  
  echo "✅ All dependencies processed"
else
  echo "ℹ️  No requirements.txt found, skipping dependencies"
fi
