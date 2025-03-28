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
PROTO_FILE="$PROTO_DIR/${SERVICE_NAME}.proto"
GEN_DIR="$SCRIPT_DIR/gen"
OUT_DIR="$GEN_DIR/$SERVICE_NAME/pb"

# check if the .proto file exists
if [ ! -f "$PROTO_FILE" ]; then
  echo "❌ Proto file not found: $PROTO_FILE"
  exit 1
fi

echo "🔄 Generating Go code from ${SERVICE_NAME}.proto ..."

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

# run go mod tidy
cd gen
echo "🔄 Running go mod tidy ..."
go mod tidy

# all done
echo "✅ Proto files for '$SERVICE_NAME' generated in $OUT_DIR"