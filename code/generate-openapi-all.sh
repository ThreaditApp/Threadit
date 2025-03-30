#!/bin/bash

PROTO_DIR="./proto"

if [ ! -d "$PROTO_DIR" ]; then
  echo "❌ Directory $PROTO_DIR does not exist."
  exit 1
fi

# loop each proto file
for proto_path in "$PROTO_DIR"/*.proto; do
  [ -e "$proto_path" ] || { echo "❌ No .proto files found in $PROTO_DIR."; exit 1; }

  filename=$(basename "$proto_path" .proto)

  echo "🌐 Generating OpenAPI spec for $filename.proto..."
  ./generate-openapi.sh "$filename"
done

echo "✅ All OpenAPI specs generated in docs/openapi/"
