#!/bin/bash

print_usage() {
  echo "Usage: $0 <folder1> [folder2 folder3 ...]"
  echo "Example: $0 user-service auth-service payment-service"
  exit 1
}

if [ $# -eq 0 ]; then
  print_usage
fi

OUTPUT_FOLDER="protogen"
mkdir -p "$OUTPUT_FOLDER"

for folder in "$@"; do
  TARGET_FOLDER="proto/$folder"

  if [ ! -d "$TARGET_FOLDER" ]; then
    echo "Warning: Folder $TARGET_FOLDER does not exist, skipping..."
    continue
  fi

  echo "Processing $TARGET_FOLDER"

  # Process all .proto files
  find "$TARGET_FOLDER" -name "*.proto" -print0 | while IFS= read -r -d '' entry; do
    echo "Processing $entry"
    protoc --proto_path=proto \
      --proto_path=proto/plugins \
      --plugin=$(go env GOPATH)/bin/protoc-gen-go \
      --plugin=$(go env GOPATH)/bin/protoc-gen-govalidators \
      --go_out="$OUTPUT_FOLDER" --go_opt=paths=source_relative \
      --go-grpc_out="$OUTPUT_FOLDER" --go-grpc_opt=paths=source_relative \
      --govalidators_out="lang=go,paths=source_relative:$OUTPUT_FOLDER" \
      "$entry"
  done

  # Process API proto files
  find "$TARGET_FOLDER" -name "*_api.proto" -print0 | while IFS= read -r -d '' entry; do
    echo "Processing API file: $entry"
    protoc --proto_path=proto \
      --proto_path=proto/plugins \
      --plugin=$(go env GOPATH)/bin/protoc-gen-grpc-gateway \
      --plugin=$(go env GOPATH)/bin/protoc-gen-openapiv2 \
      --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
      --go-grpc_out="$OUTPUT_FOLDER" --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out="$OUTPUT_FOLDER" \
      --grpc-gateway_opt=paths=source_relative,allow_delete_body=true,repeated_path_param_separator=ssv \
      --openapiv2_out="$OUTPUT_FOLDER" \
      --openapiv2_opt=logtostderr=true,repeated_path_param_separator=ssv \
      "$entry"
  done

  # Process DB proto files
  find "$TARGET_FOLDER" -name "*_db.proto" -print0 | while IFS= read -r -d '' entry; do
    echo "Processing DB file: $entry"
    protoc --proto_path=proto \
      --proto_path=proto/plugins \
      --plugin=$(go env GOPATH)/bin/protoc-gen-gorm \
      --go-grpc_out="$OUTPUT_FOLDER" --go-grpc_opt=paths=source_relative \
      --gorm_out=paths=source_relative:"$OUTPUT_FOLDER" \
      "$entry"
  done

  if [ -f "./protogen/rnd-service/rnd_api.swagger.json" ]; then
    echo "Moving swagger file from rnd-service"
    mkdir -p "./www"
    mv ./protogen/rnd-service/rnd_api.swagger.json ./www/swagger.json
  fi
done

echo "Generate complete"