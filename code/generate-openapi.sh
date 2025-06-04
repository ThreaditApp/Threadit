#!/bin/bash

# fail on any error
set -e

# Function to generate OpenAPI spec for a specific service
generate_openapi_for_service() {
  SERVICE_NAME="$1"
  SKIP_SERVICE_NAMES=("db-service" "models")

  # get paths
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

  if [[ " ${SKIP_SERVICE_NAMES[@]} " =~ " ${SERVICE_NAME} " ]]; then
    echo "✅ Skipping OpenAPI generation for ${SERVICE_NAME}"
    return
  fi

  PROTO_DIR="$SCRIPT_DIR/proto"
  GOOGLE_API_DIR="$PROTO_DIR/google/api"
  PROTO_FILE="$PROTO_DIR/${SERVICE_NAME}.proto"
  OUT_DIR="$SCRIPT_DIR/../docs/openapi"

  # check if the .proto file exists
  if [ ! -f "$PROTO_FILE" ]; then
    echo "❌ Proto file not found: $PROTO_FILE"
    return
  fi

  echo "🌐 Generating OpenAPI spec for ${SERVICE_NAME}.proto ..."

  # create output directory if it doesn't exist
  mkdir -p "$OUT_DIR"

  # run protoc with grpc-gateway's OpenAPI plugin (Swagger 2.0)
  protoc \
    --openapiv2_out="$OUT_DIR" \
    --openapiv2_opt logtostderr=true \
    --proto_path="$GOOGLE_API_DIR" \
    --proto_path="$PROTO_DIR" \
    "$PROTO_FILE"

  # Convert Swagger 2.0 JSON to OpenAPI 3.1.0 YAML using swagger2openapi
  swagger2openapi -o "$OUT_DIR/${SERVICE_NAME}.yaml" "$OUT_DIR/${SERVICE_NAME}.swagger.json"

  # Remove the JSON file
  rm -f "$OUT_DIR/${SERVICE_NAME}.swagger.json"

  echo "✅ OpenAPI spec generated at docs/openapi/gen/${SERVICE_NAME}.yaml"
}

# Check if the -s flag is provided
if [ "$1" == "-s" ] && [ -n "$2" ]; then
  # Generate for a specific service
  SERVICE_NAME="$2"
  generate_openapi_for_service "$SERVICE_NAME"

elif [ -z "$1" ]; then
  # Generate for all services
  PROTO_DIR="./proto"

  if [ ! -d "$PROTO_DIR" ]; then
    echo "❌ Directory $PROTO_DIR does not exist."
    exit 1
  fi

  # loop through each proto file
  for proto_path in "$PROTO_DIR"/*.proto; do
    [ -e "$proto_path" ] || { echo "❌ No .proto files found in $PROTO_DIR."; exit 1; }

    filename=$(basename "$proto_path" .proto)

    #echo "🌐 Generating OpenAPI spec for $filename.proto..."
    generate_openapi_for_service "$filename"
  done

  echo "✅ All OpenAPI specs generated in docs/openapi/gen"

else
  echo "❌ Invalid usage."
  echo "Usage: ./generate-openapi.sh                   # Generates for all services"
  echo "Usage: ./generate-openapi.sh -s <service-name> # Generates for a specific service"
  exit 1
fi
