#!/bin/bash

# fail on any error
set -e

# check for a service name argument
if [ -z "$1" ]; then
  echo "❌ Usage: ./generate-openapi.sh <service-name>"
  exit 1
fi

SERVICE_NAME="$1"
SKIP_SERVICE_NAMES=("db-service" "models")

# get paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ " ${SKIP_SERVICE_NAMES[@]} " =~ " ${SERVICE_NAME} " ]]; then
  echo "✅ Skipping OpenAPI generation for ${SERVICE_NAME}"
  exit 0
fi

PROTO_DIR="$SCRIPT_DIR/proto"
GOOGLE_API_DIR="$PROTO_DIR/google/api"
PROTO_FILE="$PROTO_DIR/${SERVICE_NAME}.proto"
OUT_DIR="$SCRIPT_DIR/../docs/gen-openapi"

# check if the .proto file exists
if [ ! -f "$PROTO_FILE" ]; then
  echo "❌ Proto file not found: $PROTO_FILE"
  exit 1
fi

echo "🌐 Generating OpenAPI spec for ${SERVICE_NAME}.proto ..."

# create output directory if it doesn't exist
mkdir -p "$OUT_DIR"

# clean up existing generated OpenAPI file
rm -f "$OUT_DIR/${SERVICE_NAME}.swagger.json"

# run protoc with grpc-gateway's OpenAPI plugin
protoc \
  --openapiv2_out="$OUT_DIR" \
  --openapiv2_opt logtostderr=true \
  --proto_path="$GOOGLE_API_DIR" \
  --proto_path="$PROTO_DIR" \
  "$PROTO_FILE"

# rename output file for clarity
mv "$OUT_DIR/${SERVICE_NAME}.swagger.json" "$OUT_DIR/${SERVICE_NAME}.openapi.json"

echo "✅ OpenAPI spec generated at docs/openapi/${SERVICE_NAME}.openapi.json"
