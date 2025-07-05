#!/bin/bash

set -e

OUT_DIR=$1

if [ -z "$OUT_DIR" ]; then
  echo "Usage: ./gen.sh <output-directory>"
  exit 1
fi

echo "Generating to: $OUT_DIR"

# Replace {{OUT_DIR}} in template with real path
sed "s|{{OUT_DIR}}|$OUT_DIR|g" buf.gen.yaml.template > buf.gen.yaml

# Run buf generate
buf generate --path="./proto/$OUT_DIR"

# Optional: remove temporary buf.gen.yaml
rm -f buf.gen.yaml

echo "Done ✅"