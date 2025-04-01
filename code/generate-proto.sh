#!/bin/bash

# fail on any error
set -e

# get paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="$SCRIPT_DIR/proto"
GOOGLE_API_DIR="$PROTO_DIR/google/api"
GEN_DIR="$SCRIPT_DIR/gen"

# Function to generate proto code for a single service
generate_proto() {
  SERVICE_NAME="$1"
  PROTO_FILE="$PROTO_DIR/${SERVICE_NAME}.proto"
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
    --grpc-gateway_out="." \
    --proto_path="$GOOGLE_API_DIR" \
    --proto_path="$PROTO_DIR" \
    "$PROTO_FILE"

  echo "✅ Proto files for '$SERVICE_NAME' generated in $OUT_DIR"
}

# If no argument is given, generate for all services
if [ "$#" -eq 0 ]; then
  if [ ! -d "$PROTO_DIR" ]; then
    echo "❌ Directory $PROTO_DIR does not exist."
    exit 1
  fi

  # loop through all .proto files
  for proto_path in "$PROTO_DIR"/*.proto; do
    [ -e "$proto_path" ] || { echo "❌ No .proto files found in $PROTO_DIR."; exit 1; }

    filename=$(basename "$proto_path" .proto)
    generate_proto "$filename"
  done

  echo "✅ All proto files processed."
else
  # Check for -s <service-name> flag
  if [ "$1" == "-s" ] && [ -n "$2" ]; then
    generate_proto "$2"
  else
    echo "❌ Invalid usage."
    echo "Usage:"
    echo "  ./generate-proto.sh                    # Generates for all services"
    echo "  ./generate-proto.sh -s <service-name>  # Generates for a specific service"
    exit 1
  fi
fi

# Run go mod tidy
cd "$GEN_DIR"
echo "🔄 Running go mod tidy ..."
go mod tidy

echo "✅ Done."