#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <folder>"
  exit 1
fi

TARGET_FOLDER="$1"

if [ ! -d "$TARGET_FOLDER" ]; then
  echo "Folder $TARGET_FOLDER does not exist."
  exit 1
fi

mkdir verify

echo "Processing $TARGET_FOLDER"

for entry in $(find "$TARGET_FOLDER" -name "*.proto"); do
  echo "Processing $entry"
  protoc --proto_path=. \
    --proto_path=plugins \
    --plugin=$(go env GOPATH)/bin/protoc-gen-go \
    --plugin=$(go env GOPATH)/bin/protoc-gen-govalidators \
    --go_out=verify --go_opt=paths=source_relative \
    --go-grpc_out=verify --go-grpc_opt=paths=source_relative \
    --govalidators_out="lang=go,paths=source_relative:verify" \
    $entry
done

for entry in $(find "$TARGET_FOLDER" -name "*_api.proto"); do
  echo "Processing $entry"
  protoc --proto_path=. \
    --proto_path=plugins \
    --plugin=$(go env GOPATH)/bin/protoc-gen-grpc-gateway \
    --plugin=$(go env GOPATH)/bin/protoc-gen-openapiv2 \
    --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
    --go-grpc_out=verify --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=verify \
    --grpc-gateway_opt=paths=source_relative,allow_delete_body=true,repeated_path_param_separator=ssv \
    --openapiv2_out=verify \
    --openapiv2_opt=logtostderr=true,repeated_path_param_separator=ssv \
    $entry
done

for entry in $(find "$TARGET_FOLDER" -name "*_db.proto"); do
  echo "Processing $entry"
  protoc --proto_path=. \
    --proto_path=plugins \
    --plugin=$(go env GOPATH)/bin/protoc-gen-gorm \
    --go-grpc_out=verify --go-grpc_opt=paths=source_relative \
    --gorm_out=paths=source_relative:verify \
    $entry
done

rm -rf "verify"

echo "Verify successful"